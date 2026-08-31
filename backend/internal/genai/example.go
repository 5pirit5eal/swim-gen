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
	logger.Debug("Generating prompt example...")

	audienceInstruction := ""
	if req.Audience != nil {
		switch *req.Audience {
		case models.AudienceBeginner:
			audienceInstruction = "Der Schwimmer ist ein Anfänger / lernt schwimmen. Die Anfrage soll sich auf einfache, machbare Einheiten, Wasserlage, Atmung, Technik und ausreichende Pausen konzentrieren."
		case models.AudienceTriathlete:
			audienceInstruction = "Der Schwimmer ist ein Triathlet. Die Anfrage soll sich auf Kraul-Ausdauer, Tempohärte, Freiwasser-Fokus oder aerobe Schwellenserien konzentrieren."
		case models.AudienceCompetitiveSwimmer:
			audienceInstruction = "Der Schwimmer ist ein erfahrener Leistungsschwimmer. Die Anfrage soll für längere Trainingseinheiten gedacht sein, jedoch nicht immer auf maximale Intensität abzielen. Wähle für die Einheit einen spezifischen, komplexen Trainingsschwerpunkt aus: Grundlagenausdauer (Endurance), Regeneration/Kompensation, Schwellentraining (Threshold), wettkampfspezifische Schnelligkeitsausdauer (Speed) oder Sprintschnelligkeit (Sprint)."
		case models.AudienceHobby:
			audienceInstruction = "Der Schwimmer ist ein Hobbyschwimmer / Fitness-Schwimmer. Die Anfrage soll sich auf abwechslungsreiche, motivierende Trainingseinheiten für Fitness, Spaß und Technik konzentrieren."
		}
	}

	prompt := fmt.Sprintf(generatePromptTemplateStr, audienceInstruction, req.Language)

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

	logger.Debug("Prompt generated successfully")
	return answer.Text(), nil
}
