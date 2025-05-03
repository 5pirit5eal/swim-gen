package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
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
func (db *RAGDB) Query(ctx context.Context, query string, filter map[string]any, method string) (*models.RAGResponse, error) {
	logger := httplog.LogEntry(ctx)
	// Set the embedder to query mode
	db.Client.QueryMode()
	var docs []schema.Document
	var err error
	switch {
	case query == "" && filter == nil:
		return nil, fmt.Errorf("either a query or a filter must be provided")
	case query == "" && filter != nil:
		docs, err = db.Store.Search(ctx, 10, vectorstores.WithFilters(filter))
	case query != "" && filter == nil:
		docs, err = db.Store.SimilaritySearch(ctx, query, 10)
	case query != "" && filter != nil:
		docs, err = db.Store.SimilaritySearch(ctx, query, 10, vectorstores.WithFilters(filter))
	}
	if err != nil {
		logger.Error("Error searching for documents", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error searching for documents: %w", err)
	}
	logger.Info("Documents found", "count", len(docs))
	logger.Debug("Documents:", "docs", docs)
	answer := &models.RAGResponse{}
	if method == "generate" {
		answer, err = db.Client.GeneratePlan(ctx, query, docs)
	} else if method == "choose" {
		planID, err := db.Client.ChoosePlan(ctx, query, docs)
		if err != nil {
			logger.Error("Error choosing plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error choosing plan: %w", err)
		}
		plan, err := db.GetPlan(ctx, planID, SourceOptionPlan)
		if err != nil {
			logger.Error("Error getting plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error getting plan: %w", err)
		}
		genericPlan := plan.Plan()
		answer.Title = genericPlan.Title
		answer.Description = genericPlan.Description
		answer.Table = genericPlan.Table
	} else {
		return nil, fmt.Errorf("unsupported method: %s", method)
	}
	if err != nil {
		return nil, err
	}
	logger.Info("Answer generated successfully", "answer", answer)

	return answer, nil
}

func (db *RAGDB) GetPlan(ctx context.Context, planID string, source SourceOption) (models.Planable, error) {
	logger := httplog.LogEntry(ctx)
	// Query the database for the plan with the given ID
	switch source {
	case SourceOptionPlan:
		var plan models.Plan
		err := pgxscan.Get(ctx, db.Conn, &plan, fmt.Sprintf(`SELECT plan_id, title, description, plan_table FROM %s WHERE plan_id = $1`, PlanTableName), planID)
		if err != nil {
			logger.Error("Error querying plan", httplog.ErrAttr(err))
			return nil, err
		}
		return &plan, nil
	case SourceOptionScraped:
		return db.GetScrapedPlan(ctx, planID)
	case SourceOptionDonated:
		return db.GetDonatedPlan(ctx, planID)
	}
	return nil, fmt.Errorf("unsupported source option: %s", source)
}
