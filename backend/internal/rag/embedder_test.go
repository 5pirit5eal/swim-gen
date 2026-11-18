package rag

import (
	"context"
	"testing"
)

func TestStaticEmbedderReturnsSharedVector(t *testing.T) {
	vector := []float32{0.1, 0.2, 0.3}
	embedder := staticEmbedder{vector: vector}

	got, err := embedder.EmbedQuery(context.Background(), "any query")
	if err != nil {
		t.Fatalf("EmbedQuery() error = %v", err)
	}
	if len(got) != len(vector) {
		t.Fatalf("EmbedQuery() returned %d values, want %d", len(got), len(vector))
	}
	for i := range vector {
		if got[i] != vector[i] {
			t.Errorf("EmbedQuery()[%d] = %f, want %f", i, got[i], vector[i])
		}
	}

	docs, err := embedder.EmbedDocuments(context.Background(), []string{"a", "b", "c"})
	if err != nil {
		t.Fatalf("EmbedDocuments() error = %v", err)
	}
	if len(docs) != 3 {
		t.Fatalf("EmbedDocuments() returned %d vectors, want 3", len(docs))
	}
	for i, doc := range docs {
		if len(doc) != len(vector) {
			t.Errorf("EmbedDocuments()[%d] has %d values, want %d", i, len(doc), len(vector))
		}
	}
}
