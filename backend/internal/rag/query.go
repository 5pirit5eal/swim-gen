package rag

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type SourceOption string

const (
	// SourceOptionPlan indicates that the source is a plan.
	SourceOptionPlan SourceOption = "plan"
	// SourceOptionScraped indicates that the source is a scraped plan.
	SourceOptionScraped SourceOption = "scraped"
)

// Query searches for documents in the database based on the provided query and filter.
func (db *RAGDB) Query(ctx context.Context, query string, lang models.Language, userProfile string, filter map[string]any, method string, poolLength any) (*models.Plan, error) {
	logger := httplog.LogEntry(ctx)
	// Set the embedder to query mode
	db.Client.QueryMode()
	searchQuery := buildSearchQuery(query, userProfile)
	if searchQuery == "" && filter == nil {
		return nil, fmt.Errorf("either a query or a filter must be provided")
	}
	embedder, err := db.queryEmbedder(ctx, searchQuery)
	if err != nil {
		logger.Error("Error embedding search query", httplog.ErrAttr(err))
		return nil, err
	}

	var planDocs, drillDocs []schema.Document
	var planErr, drillErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		planDocs, planErr = db.searchPlanDocuments(ctx, searchQuery, filter, embedder)
	}()
	if method == "generate" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			drillDocs, drillErr = db.searchDrillDocuments(ctx, searchQuery, lang, embedder)
		}()
	}
	wg.Wait()
	if planErr != nil {
		logger.Error("Error searching for plan documents", httplog.ErrAttr(planErr))
		return nil, fmt.Errorf("error searching for plan documents: %w", planErr)
	}
	if drillErr != nil {
		logger.Error("Error searching for drill documents", httplog.ErrAttr(drillErr))
		return nil, fmt.Errorf("error searching for drill documents: %w", drillErr)
	}

	logger.Debug("Documents found", "count", len(planDocs))
	logger.Debug("Documents:", "docs", planDocs)
	var plan models.Planable
	switch method {
	case "generate":
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

func (db *RAGDB) searchPlanDocuments(ctx context.Context, searchQuery string, filter map[string]any, embedder embeddings.Embedder) ([]schema.Document, error) {
	if searchQuery == "" {
		return db.PlanStore.Search(ctx, 5, vectorstores.WithFilters(filter))
	}
	opts := make([]vectorstores.Option, 0, 2)
	if embedder != nil {
		opts = append(opts, vectorstores.WithEmbedder(embedder))
	}
	if filter != nil {
		opts = append(opts, vectorstores.WithFilters(filter))
	}
	return db.PlanStore.SimilaritySearch(ctx, searchQuery, 5, opts...)
}

func (db *RAGDB) searchDrillDocuments(ctx context.Context, searchQuery string, lang models.Language, embedder embeddings.Embedder) ([]schema.Document, error) {
	opts := []vectorstores.Option{vectorstores.WithFilters(map[string]any{"language": string(lang)})}
	if embedder != nil {
		opts = append(opts, vectorstores.WithEmbedder(embedder))
	}
	return db.DrillStore.SimilaritySearch(ctx, searchQuery, 10, opts...)
}

func buildSearchQuery(query, userProfile string) string {
	return strings.TrimSpace(strings.Join([]string{
		strings.TrimSpace(query),
		strings.TrimSpace(userProfile),
	}, "\n\n"))
}
