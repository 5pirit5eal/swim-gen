package genai

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Creates Embeddings according to the Langchaingo interface with the Google GenAI package
func (c *GoogleGenAIClient) CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error) {
	contents := make([]*genai.Content, len(texts))
	for i, text := range texts {
		contents[i] = genai.NewContentFromText(text, genai.RoleUser)
	}

	resp, err := c.gc.Models.EmbedContent(ctx, c.cfg.Embedding.Model, contents, c.embedCfg)
	if err != nil {
		return nil, fmt.Errorf("Models.EmbedContent: %w", err)
	}
	embeddings := make([][]float32, len(resp.Embeddings))
	for i, embedding := range resp.Embeddings {
		embeddings[i] = embedding.Values
	}
	return embeddings, nil
}

func (c *GoogleGenAIClient) QueryMode() {
	c.embedCfg.TaskType = "RETRIEVAL_QUERY"
}

func (c *GoogleGenAIClient) DocumentMode() {
	c.embedCfg.TaskType = "RETRIEVAL_DOCUMENT"
}
