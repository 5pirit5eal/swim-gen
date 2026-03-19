package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"google.golang.org/genai"
)

func (gc *GoogleGenAIClient) ImprovePlan(ctx context.Context, plan models.ScrapedPlan, syncGroup *sync.WaitGroup, c chan<- models.Document, ec chan<- error) {
	if syncGroup != nil {
		defer syncGroup.Done()
	}
	logger := slog.Default()

	// Step 0: Sanitize input to prevent encoding issues and data pollution
	plan.Title = models.SanitizeString(plan.Title)
	plan.Description = models.SanitizeString(plan.Description)
	models.SanitizeRows(&plan.Table)

	// Step 1: Restructure plan with nested loops support (happens only once during scraping)
	logger.Debug("Restructuring plan", "plan_title", plan.Plan().Title)
	restructuredPlan, err := gc.RestructurePlan(ctx, plan.Plan())
	if err != nil {
		logger.Warn("Failed to restructure plan, using original", "plan_title", plan.Plan().Title, httplog.ErrAttr(err))
		restructuredPlan = plan.Plan()
	}

	// Update the original plan with the restructured version (this preserves the Planable type and URL)
	plan.Title = restructuredPlan.Title
	plan.Description = restructuredPlan.Description
	plan.Table = restructuredPlan.Table

	// Step 2: Generate metadata
	logger.Debug("Generating metadata", "plan_title", restructuredPlan.Title)
	meta, err := gc.GenerateMetadata(ctx, restructuredPlan)
	if err != nil {
		logger.Error("Error when generating metadata with LLM", httplog.ErrAttr(err))
		ec <- fmt.Errorf("error generating metadata: %w", err)
		return
	}

	// Step 3: Create document with restructured plan
	c <- models.Document{Plan: &plan, Meta: meta}
}

// RestructurePlan analyzes a plan and optimizes its structure by identifying repeating patterns
// and representing them using nested SubRows instead of flat rows. This happens only once during
// the scraping/import process. The actual training content is NEVER modified - only the schema
// representation is optimized for better structure.
func (gc *GoogleGenAIClient) RestructurePlan(ctx context.Context, plan *models.Plan) (*models.Plan, error) {
	logger := slog.Default()
	gps, err := models.GeneratedPlanSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get GeneratedPlan schema: %w", err)
	}

	// Convert table to json
	genericPlan := plan.Plan()
	genericPlan.Table = annotateTable(genericPlan.Table)
	query := fmt.Sprintf(restructureTemplateStr, gps, genericPlan.Title, genericPlan.Description, genericPlan.Table.String())

	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	genCfg.ResponseJsonSchema = gps
	logger.Debug("Requesting restructuring by LLM", "plan", genericPlan.String())
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)

	if err != nil {
		logger.Error("Error when restructuring plan", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error restructuring plan: %w", err)
	}

	var gp models.GeneratedPlan
	err = json.Unmarshal([]byte(answer.Text()), &gp)
	if err != nil {
		logger.Error("Error parsing restructured plan", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing restructured plan: %w", err)
	}
	logger.Debug("Restructured plan", "plan", gp.Map())

	gp.Table.UpdateSum()
	if err := gp.Table.Validate(); err != nil {
		logger.Error("Validation failed for restructured plan", httplog.ErrAttr(err))
		return nil, fmt.Errorf("validation failed for restructured plan: %w", err)
	}

	// This preserves the original Planable type (e.g., ScrapedPlan with URL)
	genericPlan.Title = gp.Title
	genericPlan.Description = gp.Description
	genericPlan.Table = gp.Table

	// Return the original Planable interface (preserves URL and other fields)
	return genericPlan, nil
}

// annotateTable annotates each row with a note if it contains a "+" or " x " character,
// which are common indicators of potential subrows in swim training plans.
// This helps the LLM identify patterns that may indicate nested structures during restructuring.
func annotateTable(table models.Table) models.Table {
	for i := range table {
		if strings.Contains(table[i].Content, "+") {
			table[i].Content += " [→ SUBROW-KANDIDAT: '+' Zeichen gefunden]"
		}

		if strings.Contains(table[i].Content, " x ") {
			table[i].Content += " [→ SUBROW-KANDIDAT: ' x ' Zeichen gefunden]"
		}
	}
	return table
}

func (gc *GoogleGenAIClient) GenerateMetadata(ctx context.Context, plan *models.Plan) (*models.Metadata, error) {
	logger := slog.Default()
	ms, err := models.MetadataSchema()
	if err != nil {
		logger.Error("Failed in retrieving Schema", httplog.ErrAttr(err))
		return nil, fmt.Errorf("models.MetadataSchema: %w", err)
	}

	// Convert table to json
	genericPlan := plan.Plan()
	tableJSON, err := json.Marshal(genericPlan.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal table to JSON: %w", err)
	}
	// Enhance scraped documents with gemini and create meaningful metadata
	query := fmt.Sprintf(metadataTemplateStr, models.Abbreviations, genericPlan.Title, genericPlan.Description, string(tableJSON), ms)
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("Models.GenerateContent: %w", err)
	}
	logger.Debug("Successful answer from LLM", "answer", answer.Text())

	// Parse the answer as JSON
	var metadata models.Metadata
	err = json.Unmarshal([]byte(answer.Text()), &metadata)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("JSON unmarshal error: %w with raw response %s", err, answer.Text())
	}

	return &metadata, nil
}
