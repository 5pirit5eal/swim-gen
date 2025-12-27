package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type SourceOption string

const (
	// SourceOptionPlan indicates that the source is a plan.
	SourceOptionPlan SourceOption = "plan"
	// SourceOptionScraped indicates that the source is a scraped plan.
	SourceOptionScraped SourceOption = "scraped"
	// SourceOptionDonated indicates that the source is a donated plan.
	SourceOptionDonated SourceOption = "donated"
)

// Query searches for documents in the database based on the provided query and filter.
func (db *RAGDB) Query(ctx context.Context, query string, lang models.Language, userProfile string, filter map[string]any, method string, poolLength any) (*models.Plan, error) {
	logger := httplog.LogEntry(ctx)
	// Set the embedder to query mode
	db.Client.QueryMode()
	var planDocs []schema.Document
	var err error
	switch {
	case query == "" && filter == nil:
		return nil, fmt.Errorf("either a query or a filter must be provided")
	case query == "" && filter != nil:
		planDocs, err = db.PlanStore.Search(ctx, 5, vectorstores.WithFilters(filter))
	case query != "" && filter == nil:
		planDocs, err = db.PlanStore.SimilaritySearch(ctx, query, 5)
	case query != "" && filter != nil:
		planDocs, err = db.PlanStore.SimilaritySearch(ctx, query, 5, vectorstores.WithFilters(filter))
	}
	if err != nil {
		logger.Error("Error searching for plan documents", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error searching for plan documents: %w", err)
	}

	logger.Debug("Documents found", "count", len(planDocs))
	logger.Debug("Documents:", "docs", planDocs)
	var plan models.Planable
	switch method {
	case "generate":
		var drillDocs []schema.Document
		drillDocs, err = db.DrillStore.SimilaritySearch(ctx, query, 10, vectorstores.WithFilters(map[string]any{"language": string(lang)}))
		if err != nil {
			logger.Error("Error searching for drill documents", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error searching for drill documents: %w", err)
		}

		plan, err = db.Client.GeneratePlan(ctx, query, string(lang), userProfile, poolLength, planDocs, drillDocs)
		if err != nil {
			logger.Error("Error generating plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error generating plan: %w", err)
		}
	case "choose":
		if len(planDocs) == 0 {
			return nil, fmt.Errorf("no documents in database matching query and filters")
		}
		var planID string
		planID, err = db.Client.ChoosePlan(ctx, query, string(lang), poolLength, planDocs)
		if err != nil {
			logger.Error("Error choosing plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error choosing plan: %w", err)
		}

		plan, err = db.GetPlan(ctx, planID, SourceOptionPlan)
		if err != nil {
			logger.Error("Error getting plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error getting plan: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	genericPlan := plan.Plan()

	if lang != "de" {
		genericPlan, err = db.Client.TranslatePlan(ctx, genericPlan, lang)
		if err != nil {
			logger.Error("Error translating plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error translating plan: %w", err)
		}

	}

	logger.Debug("Plan generated successfully", "plan", genericPlan)

	return genericPlan, nil
}
