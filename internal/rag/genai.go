package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

// GeneratePlan generates a plan using the LLM based on the provided query and documents.
func (db *RAGDB) GeneratePlan(ctx context.Context, q string, docs []schema.Document) (*models.RAGResponse, error) {
	logger := httplog.LogEntry(ctx)
	ts, err := models.TableSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get table schema: %w", err)
	}

	var dc []string
	for _, doc := range docs {
		dc = append(dc, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(ragTemplateStr, ts, q, strings.Join(dc, "\n \n"))
	answer, err := llms.GenerateFromSinglePrompt(ctx, db.Client, query,
		llms.WithResponseMIMEType("application/json"), llms.WithTemperature(float64(1.0)))
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error generating answer: %w", err)
	}
	answer = cleanResponse(answer)
	// read description and table from the LLM response
	var p models.RAGResponse
	err = json.Unmarshal([]byte(answer), &p)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer)
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	// Add the total to the table if it is not already present
	if !strings.Contains(p.Table[len(p.Table)-1].Content, "Total") {
		p.Table.AddSum()
	}
	// Recalculate the sums of the rows to be sure they are correct
	p.Table.UpdateSum()

	// Add the plan to the response
	logger.Debug("Plan generated successfully")
	return &p, nil
}

// ChoosePlan lets an LLM choose the best fitting plan from the given documents.
func (db *RAGDB) ChoosePlan(ctx context.Context, q string, docs []schema.Document) (*models.RAGResponse, error) {
	logger := httplog.LogEntry(ctx)
	var dc string
	for i, doc := range docs {
		dc += fmt.Sprintf("%d: %s \n\n", i, doc.PageContent)
	}

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(choosePlanTemplateStr, q, dc)
	answer, err := llms.GenerateFromSinglePrompt(ctx, db.Client, query, llms.WithResponseMIMEType("application/json"))
	if err != nil {
		logger.Error("Error when generating answer with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error generating answer: %w", err)
	}
	log.Println("Answer from LLM:", answer)
	answer = cleanResponse(answer)
	var cr models.ChooseResponse
	err = json.Unmarshal([]byte(answer), &cr)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer)
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	var t models.Table
	docString, err := json.Marshal(docs[cr.Idx].Metadata["table"])
	logger.Info("Generated docstring", "table", string(docString))
	if err != nil {
		logger.Error("Error parsing table from LLM response, json.Marshal:", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error parsing table from LLM response: %w", err)
	}
	err = json.Unmarshal(docString, &t)
	if err != nil {
		logger.Error("Error parsing table from LLM response, json.Unmarshal:", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error parsing table from LLM response: %w", err)
	}

	return &models.RAGResponse{
		Title:       docs[cr.Idx].Metadata["title"].(string),
		Description: "Begründung:" + cr.Description + "\n" + docs[cr.Idx].Metadata["description"].(string),
		Table:       t,
	}, nil
}

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
		answer, err = db.GeneratePlan(ctx, query, docs)
	} else {
		query = fmt.Sprintf("Ich suche nach einem Plan mit folgenden Kriterien: %v", filter)
		answer, err = db.ChoosePlan(ctx, query, docs)
	}
	if err != nil {
		return nil, err
	}
	logger.Info("Answer generated successfully", "answer", answer)

	return answer, nil
}

func cleanResponse(s string) string {
	s, _ = strings.CutPrefix(s, "```json")
	s, _ = strings.CutSuffix(s, "```")
	return s
}
