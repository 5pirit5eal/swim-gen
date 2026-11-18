package rag

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/embeddings"
)

type staticEmbedder struct {
	vector []float32
}

var _ embeddings.Embedder = staticEmbedder{}

func (s staticEmbedder) EmbedQuery(context.Context, string) ([]float32, error) {
	return s.vector, nil
}

func (s staticEmbedder) EmbedDocuments(_ context.Context, texts []string) ([][]float32, error) {
	vectors := make([][]float32, len(texts))
	for i := range texts {
		vectors[i] = s.vector
	}
	return vectors, nil
}

func (db *RAGDB) queryEmbedder(ctx context.Context, searchQuery string) (embeddings.Embedder, error) {
	if searchQuery == "" {
		return nil, nil
	}
	embs, err := db.Client.CreateEmbedding(ctx, []string{searchQuery})
	if err != nil {
		return nil, fmt.Errorf("error embedding search query: %w", err)
	}
	return staticEmbedder{vector: embs[0]}, nil
}
