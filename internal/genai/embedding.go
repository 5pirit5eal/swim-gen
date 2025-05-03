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

	// Embed in batches to avoid exceeding the max tokens
	batchSize := 10
	embeddings := make([][]float32, len(texts))
	for i := 0; i < (len(texts)+batchSize-1)/batchSize; i++ {
		if i*batchSize >= len(texts) {
			break
		}
		var batch []*genai.Content
		if i*batchSize+batchSize > len(texts) {
			batch = contents[i*batchSize:]
		} else {
			batch = contents[i*batchSize : i*batchSize+batchSize]
		}
		resp, err := c.gc.Models.EmbedContent(ctx, c.cfg.Embedding.Model, batch, c.embedCfg)
		if err != nil {
			return nil, fmt.Errorf("Models.EmbedContent: %w", err)
		}
		for j, embedding := range resp.Embeddings {
			embeddings[i*batchSize+j] = embedding.Values
		}
	}

	return embeddings, nil
}

func (c *GoogleGenAIClient) QueryMode() {
	c.embedCfg.TaskType = "RETRIEVAL_QUERY"
}

func (c *GoogleGenAIClient) DocumentMode() {
	c.embedCfg.TaskType = "RETRIEVAL_DOCUMENT"
}
