package genai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
	"google.golang.org/genai"
)

// GeneratePlan generates a plan using the LLM based on the provided query and documents.
func (gc *GoogleGenAIClient) GeneratePlan(ctx context.Context, q, lang string, poolLength any, docs []schema.Document) (*models.GeneratedPlan, error) {
	logger := httplog.LogEntry(ctx)
	gps, err := models.GeneratedPlanSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get GeneratedPlan schema: %w", err)
	}

	var dc []string
	for _, doc := range docs {
		dc = append(dc, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(ragTemplateStr, poolLength, lang, q, strings.Join(dc, "\n \n"))
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	genCfg.ResponseJsonSchema = gps
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)

	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error when generating answer with LLM: %w", err)
	}

	// read description and table from the LLM response
	var p models.GeneratedPlan
	err = json.Unmarshal([]byte(answer.Text()), &p)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	// Add the total to the table if it is not already present
	if !strings.Contains(p.Table[len(p.Table)-1].Content, "Gesamt") {
		p.Table.AddSum()
	}
	// Recalculate the sums of the rows to be sure they are correct
	p.Table.UpdateSum()

	// Add the plan to the response
	logger.Debug("Plan generated successfully")
	return &p, nil
}

// ChoosePlan lets an LLM choose the best fitting plan from the given documents.
// Returns the plan id of the chosen plan
func (gc *GoogleGenAIClient) ChoosePlan(ctx context.Context, q, lang string, poolLength any, docs []schema.Document) (string, error) {
	logger := httplog.LogEntry(ctx)
	var dc string
	for i, doc := range docs {
		dc += fmt.Sprintf("%d: %s \n\n", i, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(choosePlanTemplateStr, poolLength, lang, q, dc)
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return "", fmt.Errorf("error generating answer: %w", err)
	}
	logger.Debug("Successful answer from LLM", "answer", answer)

	var cr models.ChoiceResult
	err = json.Unmarshal([]byte(answer.Text()), &cr)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer)
		return "", fmt.Errorf("error parsing LLM response: %w", err)
	}
	planID, ok := docs[cr.Idx].Metadata["plan_id"]
	if !ok {
		return "", fmt.Errorf("plan_id not found in Metadata for document at index %d", cr.Idx)
	}
	planIDStr, ok := planID.(string)
	if !ok {
		return "", fmt.Errorf("plan_id is not a string in Metadata for document at index %d", cr.Idx)
	}
	return planIDStr, nil
}

func (gc *GoogleGenAIClient) ImprovePlan(ctx context.Context, plan models.Planable, syncGroup *sync.WaitGroup, c chan<- models.Document, ec chan<- error) {
	if syncGroup != nil {
		defer syncGroup.Done()
	}
	logger := httplog.LogEntry(ctx)
	meta, err := gc.GenerateMetadata(ctx, plan)
	if err != nil {
		logger.Error("Error when generating metadata with LLM", httplog.ErrAttr(err))
		ec <- fmt.Errorf("error generating metadata: %w", err)
		return
	}

	// Create request body by converting the plans into documents
	c <- models.Document{Plan: plan, Meta: meta}
}

func (gc *GoogleGenAIClient) DescribeTable(ctx context.Context, table *models.Table) (*models.Description, error) {
	logger := httplog.LogEntry(ctx)
	ds, err := models.DescriptionSchema()
	if err != nil {
		logger.Error("Failed in retrieving Schema", httplog.ErrAttr(err))
		return nil, fmt.Errorf("models.MetadataSchema: %w", err)
	}
	// Create a description of the table
	query := fmt.Sprintf(describeTemplateStr, ds, table.String())
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		return nil, fmt.Errorf("Models.GenerateContent: %w", err)
	}
	var desc models.Description
	err = json.Unmarshal([]byte(answer.Text()), &desc)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	return &desc, nil
}

func (gc *GoogleGenAIClient) GenerateMetadata(ctx context.Context, plan models.Planable) (*models.Metadata, error) {
	logger := httplog.LogEntry(ctx)
	ms, err := models.MetadataSchema()
	if err != nil {
		logger.Error("Failed in retrieving Schema", httplog.ErrAttr(err))
		return nil, fmt.Errorf("models.MetadataSchema: %w", err)
	}
	// Enhance scraped documents with gemini and create meaningful metadata
	genericPlan := plan.Plan()
	query := fmt.Sprintf(metadataTemplateStr, models.Abbreviations, genericPlan.Title, genericPlan.Description, genericPlan.Table.String(), ms)
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

// TranslatePlan translates the given plan into the specified language.
//
// Returns a copy of the plan translated to the target language.
func (gc *GoogleGenAIClient) TranslatePlan(ctx context.Context, plan *models.Plan, lang models.Language) (*models.Plan, error) {
	logger := httplog.LogEntry(ctx)
	gps, err := models.GeneratedPlanSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get GeneratedPlan schema: %w", err)
	}

	// Translate the plan to the requested language
	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(translateTemplateStr, lang, models.Abbreviations, plan.Title, plan.Description, plan.Table.String())
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	genCfg.ResponseJsonSchema = gps
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error when generating answer with LLM: %w", err)
	}

	var gp models.GeneratedPlan
	err = json.Unmarshal([]byte(answer.Text()), &gp)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}

	p := models.Plan{
		PlanID:      plan.PlanID,
		Title:       gp.Title,
		Description: gp.Description,
		Table:       gp.Table,
	}
	logger.Debug("Plan translated successfully", "plan_id", p.PlanID)
	return &p, nil
}

func (gc *GoogleGenAIClient) FileToPlan(ctx context.Context, file []byte, filename string, mimeType string, language models.Language) (*models.GeneratedPlan, error) {
	logger := httplog.LogEntry(ctx)
	logger.Debug("FileToPlan", "filename", filename, "mimeType", mimeType)
	gps, err := models.GeneratedPlanSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get GeneratedPlan schema: %w", err)
	}

	prompt := fmt.Sprintf(ocrTemplateStr, language)

	// Create a RAG query for the LLM with the most relevant documents as context
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	genCfg.ResponseJsonSchema = gps
	parts := []*genai.Part{
		genai.NewPartFromBytes(file, mimeType),
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, contents, &genCfg)

	if err != nil {
		logger.Error("Error when extracting plan from image with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error when extracting plan from image with LLM: %w", err)
	}

	// read description and table from the LLM response
	var p models.GeneratedPlan
	err = json.Unmarshal([]byte(answer.Text()), &p)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	// Add the total to the table if it is not already present
	if len(p.Table) == 0 || !strings.Contains(p.Table[len(p.Table)-1].Content, "Gesamt") {
		p.Table.AddSum()
	}
	// Recalculate the sums of the rows to be sure they are correct
	p.Table.UpdateSum()

	// Add the plan to the response
	logger.Debug("Plan extracted from image successfully")
	return &p, nil

}
