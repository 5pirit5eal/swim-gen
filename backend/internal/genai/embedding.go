package genai

import (
	"context"
	"fmt"

	"github.com/go-chi/httplog/v2"
	"google.golang.org/genai"
)

// Creates Embeddings according to the Langchaingo interface with the Google GenAI package
func (c *GoogleGenAIClient) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	logger := httplog.LogEntry(ctx)
	embeddings := make([][]float32, len(texts))
	// Gemini Embedding 2's embedContent endpoint accepts one content per request.
	for i, text := range texts {
		input := text
		if c.queryMode {
			input = formatQueryEmbeddingInput(text)
		}
		content := genai.NewContentFromText(input, genai.RoleUser)
		finishInference := logInference(logger, "create_embedding", c.cfg.Embedding.Model)
		resp, err := c.gc.Models.EmbedContent(ctx, c.cfg.Embedding.Model, []*genai.Content{content}, c.embedCfg)
		finishInference(err)
		if err != nil {
			return nil, fmt.Errorf("Models.EmbedContent: %w", err)
		}
		if len(resp.Embeddings) != 1 {
			return nil, fmt.Errorf("Models.EmbedContent: expected one embedding, got %d", len(resp.Embeddings))
		}
		embeddings[i] = resp.Embeddings[0].Values
	}

	return embeddings, nil
}

func formatQueryEmbeddingInput(content string) string {
	return fmt.Sprintf("task: search result | query: %s", content)
}

func (c *GoogleGenAIClient) QueryMode() {
	c.queryMode = true
}

func (c *GoogleGenAIClient) DocumentMode() {
	c.queryMode = false
}
