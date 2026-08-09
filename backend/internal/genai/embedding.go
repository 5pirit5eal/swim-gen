package genai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Creates Embeddings according to the Langchaingo interface with the Google GenAI package
func (c *GoogleGenAIClient) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))
	// Gemini Embedding 2's embedContent endpoint accepts one content per request.
	for i, text := range texts {
		input := text
		if c.queryMode {
			input = formatQueryEmbeddingInput(text)
		}
		content := genai.NewContentFromText(input, genai.RoleUser)
		resp, err := c.gc.Models.EmbedContent(ctx, c.cfg.Embedding.Model, []*genai.Content{content}, c.embedCfg)
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
