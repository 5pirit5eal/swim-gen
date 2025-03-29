package rag

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

func Answer(ctx context.Context, llm llms.Model, q string, docs []schema.Document) (string, error) {
	var dc []string
	for _, doc := range docs {
		dc = append(dc, doc.PageContent)
	}

	log.Printf("Found %d documents", len(dc))
	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(ragTemplateStr, q, strings.Join(dc, "\n"))
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		return "", fmt.Errorf("error generating answer: %w", err)
	}
	return answer, nil
}
