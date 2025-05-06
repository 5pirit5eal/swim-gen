package genai_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/5pirit5eal/swim-rag/internal/config"
	"github.com/5pirit5eal/swim-rag/internal/genai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEmbedding(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadConfig("../../.env", true)
	require.NoError(t, err)
	client, err := genai.NewGoogleGenAIClient(ctx, cfg)
	require.NoError(t, err)
	texts := []string{"Hello", "World"}

	_, err = client.CreateEmbedding(ctx, texts)
	assert.NoError(t, err)
}

func TestCreateEmbeddingWithManyTexts(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadConfig("../../.env", true)
	require.NoError(t, err)
	client, err := genai.NewGoogleGenAIClient(ctx, cfg)
	require.NoError(t, err)

	texts := make([]string, 190)
	for i := 0; i < 190; i++ {
		texts[i] = fmt.Sprintf("Text %d", i+1)
	}

	embeddings, err := client.CreateEmbedding(ctx, texts)
	assert.NoError(t, err)
	assert.Equal(t, 190, len(embeddings))
}
