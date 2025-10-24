package genai

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"google.golang.org/genai"
)

func (c *GoogleGenAIClient) GeneratePrompt(ctx context.Context, req models.GeneratePromptRequest) (string, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("Generating prompt example...")

	prompt := fmt.Sprintf(generatePromptTemplateStr, req.Language)

	gcfg := &genai.GenerateContentConfig{
		CandidateCount: int32(1),
		Temperature:    genai.Ptr(float32(1.5)),
		SafetySettings: []*genai.SafetySetting{
			{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
			{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
			{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
			{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
		},
	}
	answer, err := c.gc.Models.GenerateContent(ctx, c.cfg.SmallModel, genai.Text(prompt), gcfg)
	if err != nil {
		logger.Error("Error generating answer", httplog.ErrAttr(err))
		return "", err
	}

	logger.Info("Prompt generated successfully")
	return answer.Text(), nil
}
