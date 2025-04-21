package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func (s *Scraper) ImprovePlan(ctx context.Context, plan models.Plan, c chan schema.Document, ec chan error) {
	logger := httplog.LogEntry(ctx)
	ms, err := models.MetadataSchema()
	if err != nil {
		logger.Error("Failed in retrieving Schema", httplog.ErrAttr(err))
		ec <- err
		return
	}
	// Enhance scraped documents with gemini and create meaningful metadata
	query := fmt.Sprintf(scrapeTemplateStr, plan.Title, plan.Description, plan.Table.String(), ms)
	answer, err := llms.GenerateFromSinglePrompt(ctx, s.db.Client, query, llms.WithResponseMIMEType("application/json"))
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		ec <- fmt.Errorf("LLM generation error: %w", err)
		return
	}
	logger.Debug("Successful answer from LLM", "answer", answer)

	planMap := plan.Map()
	// Parse the answer as JSON
	var metadata models.Metadata
	err = json.Unmarshal([]byte(answer), &metadata)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err))
		ec <- fmt.Errorf("JSON unmarshal error: %w", err)
		return
	}

	// Add the results to the map
	maps.Copy(planMap, models.StructToMap(metadata))

	// Add the description to the plan descriptions
	plan.Description += "\n" + metadata.Reasoning
	// Create request body by converting the models.plans into documents
	c <- schema.Document{
		PageContent: plan.String(),
		Metadata:    planMap,
	}
}
