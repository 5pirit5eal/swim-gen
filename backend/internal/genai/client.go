package genai

import (
	"context"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"google.golang.org/genai"
)

type GoogleGenAIClient struct {
	gc       *genai.Client
	gcfg     *genai.GenerateContentConfig
	embedCfg *genai.EmbedContentConfig
	cfg      config.Config
}

func NewGoogleGenAIClient(ctx context.Context, cfg config.Config) (*GoogleGenAIClient, error) {
	gc, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  cfg.ProjectID,
		Location: cfg.Region,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		return nil, err
	}
	gcfg := &genai.GenerateContentConfig{
		CandidateCount: int32(1),
		Temperature:    genai.Ptr(float32(1.5)),
		SafetySettings: []*genai.SafetySetting{
			{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
			{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
			{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
			{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: false,
			ThinkingBudget:  genai.Ptr[int32](0),
		},
	}
	embedCfg := &genai.EmbedContentConfig{
		// Default embedding task type
		TaskType:             "RETRIEVAL_DOCUMENT",
		OutputDimensionality: genai.Ptr(int32(cfg.Embedding.Size)),
	}
	return &GoogleGenAIClient{
		gc:       gc,
		gcfg:     gcfg,
		embedCfg: embedCfg,
		cfg:      cfg,
	}, nil
}
