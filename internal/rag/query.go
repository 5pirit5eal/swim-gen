package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// Query searches for documents in the database based on the provided query and filter.
func (db *RAGDB) Query(ctx context.Context, query string, filter map[string]any) (*models.RAGResponse, error) {
	logger := httplog.LogEntry(ctx)
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
	var answer *models.RAGResponse
	if query != "" {
		answer, err = db.Client.GeneratePlan(ctx, query, docs)
	} else {
		query = fmt.Sprintf("Ich suche nach einem Plan mit folgenden Kriterien: %v", filter)
		answer, err = db.Client.ChoosePlan(ctx, query, docs)
	}
	if err != nil {
		return nil, err
	}
	logger.Info("Answer generated successfully", "answer", answer)

	return answer, nil
}
