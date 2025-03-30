package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

func GeneratePlan(ctx context.Context, llm llms.Model, q string, docs []schema.Document) (*models.RAGResponse, error) {
	ts, err := models.TableSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get table schema: %w", err)
	}

	var dc []string
	for _, doc := range docs {
		dc = append(dc, doc.PageContent)
	}

	log.Printf("Found %d documents", len(dc))
	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(ragTemplateStr, ts, q, strings.Join(dc, "\n \n"))
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		return nil, fmt.Errorf("error generating answer: %w", err)
	}
	answer = cleanResponse(answer)
	var p models.RAGResponse
	err = json.Unmarshal([]byte(answer), &p)
	if err != nil {
		log.Println("Error parsing LLM response:", err)
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	// Add the total to the table if it is not already present
	if !strings.Contains(p.Table[len(p.Table)-1].Content, "Total") {
		p.Table.AddSum()
	}
	p.Plan = p.Table.String()
	return &p, nil
}

func ChoosePlan(ctx context.Context, llm llms.Model, q string, docs []schema.Document) (*models.RAGResponse, error) {
	var dc string
	for i, doc := range docs {
		dc += fmt.Sprintf("%d: %s \n\n", i, doc.PageContent)
	}

	log.Printf("Found %d documents", len(dc))
	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(choosePlanTemplateStr, q, dc)
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		return nil, fmt.Errorf("error generating answer: %w", err)
	}
	answer = cleanResponse(answer)
	var cr models.ChooseResponse
	err = json.Unmarshal([]byte(answer), &cr)
	if err != nil {
		log.Println("Error parsing LLM response:", err)
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}
	var t models.Table
	err = json.Unmarshal([]byte(docs[cr.Idx].Metadata["table"].(string)), &t)
	if err != nil {
		log.Println("Error parsing Table:", err)
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}

	return &models.RAGResponse{
		Description: cr.Description,
		Plan:        docs[cr.Idx].PageContent,
		Table:       t,
	}, nil
}

func (d *RAGDB) Query(ctx context.Context, client llms.Model, q string, f map[string]any) (*models.RAGResponse, error) {
	// Find the most similar documents.
	var docs []schema.Document
	var err error
	switch {
	case q == "" && f == nil:
		return nil, fmt.Errorf("either a query or a filter must be provided")
	case q == "" && f != nil:
		docs, err = d.Store.Search(ctx, 10, vectorstores.WithFilters(f))
	case q != "" && f == nil:
		docs, err = d.Store.SimilaritySearch(ctx, q, 10)
	case q != "" && f != nil:
		docs, err = d.Store.SimilaritySearch(ctx, q, 10, vectorstores.WithFilters(f))
	}
	if err != nil {
		return nil, fmt.Errorf("error searching for documents: %w", err)
	}
	var answer *models.RAGResponse
	if q != "" {
		answer, err = GeneratePlan(ctx, client, q, docs)
	} else {
		q = fmt.Sprintf("Ich suche nach einem Plan mit folgenden Kriterien: %v", f)
		answer, err = ChoosePlan(ctx, client, q, docs)
	}
	if err != nil {
		return nil, err
	}

	return answer, nil
}

func cleanResponse(s string) string {
	s, _ = strings.CutPrefix(s, "```json")
	s, _ = strings.CutSuffix(s, "```")
	return s
}
