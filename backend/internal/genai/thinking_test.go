package genai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

func TestInferenceThinkingLevelsPreserveDefaults(t *testing.T) {
	ctx := context.Background()
	plan := &models.Plan{
		Title: "Technique session", Description: "Freestyle",
		Table: models.Table{{Amount: 4, Distance: 50, Content: "Freestyle"}},
	}
	planJSON, err := json.Marshal(models.GeneratedPlan{
		Title: plan.Title, Description: plan.Description, Table: plan.Table,
	})
	require.NoError(t, err)
	chatJSON, err := json.Marshal(models.ChatResponse{Response: "Updated", Plan: nil})
	require.NoError(t, err)

	var request struct {
		GenerationConfig genai.GenerateContentConfig `json:"generationConfig"`
		SafetySettings   []*genai.SafetySetting      `json:"safetySettings"`
	}
	var requestPath, response string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestPath = r.URL.Path
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&request))
		w.Header().Set("Content-Type", "application/json")
		assert.NoError(t, json.NewEncoder(w).Encode(&genai.GenerateContentResponse{
			Candidates: []*genai.Candidate{{Content: genai.NewContentFromText(response, genai.RoleModel)}},
		}))
	}))
	defer server.Close()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendVertexAI, Project: "test", Location: "test",
		HTTPClient: server.Client(), HTTPOptions: genai.HTTPOptions{BaseURL: server.URL},
	})
	require.NoError(t, err)

	// A non-default level catches mutations through the shared ThinkingConfig pointer.
	defaults := &genai.GenerateContentConfig{
		Temperature: genai.Ptr(float32(1.5)),
		SafetySettings: []*genai.SafetySetting{
			{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockThresholdBlockLowAndAbove},
		},
		ThinkingConfig: &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelHigh},
	}
	defaultsJSON, err := json.Marshal(defaults)
	require.NoError(t, err)
	gc := &GoogleGenAIClient{
		gc: client, gcfg: defaults,
		cfg: config.Config{Model: "test-model", SmallModel: "test-small-model"},
	}

	tests := []struct {
		name     string
		level    genai.ThinkingLevel
		model    string
		response string
		invoke   func() error
	}{
		{
			name: "prompt example", level: genai.ThinkingLevelMinimal, model: gc.cfg.SmallModel,
			response: "Generate a freestyle session",
			invoke: func() error {
				_, err := gc.GeneratePrompt(ctx, models.GeneratePromptRequest{Language: models.LanguageEN})
				return err
			},
		},
		{
			name: "plan generation", level: genai.ThinkingLevelLow, model: gc.cfg.Model,
			response: string(planJSON),
			invoke: func() error {
				_, err := gc.GeneratePlan(ctx, "Technique", "en", "", 25, nil, nil)
				return err
			},
		},
		{
			name: "translation", level: genai.ThinkingLevelLow, model: gc.cfg.Model,
			response: string(planJSON),
			invoke: func() error {
				_, err := gc.TranslatePlan(ctx, plan, models.LanguageDE)
				return err
			},
		},
		{
			name: "chat refinement", level: genai.ThinkingLevelLow, model: gc.cfg.Model,
			response: string(chatJSON),
			invoke: func() error {
				_, err := gc.ChatRefine(ctx, "", plan, "Make it easier", "en", 25, nil, nil)
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request.GenerationConfig = genai.GenerateContentConfig{}
			response = tt.response
			require.NoError(t, tt.invoke())

			got := request.GenerationConfig
			require.NotNil(t, got.ThinkingConfig)
			assert.Equal(t, tt.level, got.ThinkingConfig.ThinkingLevel)
			assert.Nil(t, got.ThinkingConfig.ThinkingBudget)
			assert.False(t, got.ThinkingConfig.IncludeThoughts)
			assert.Equal(t, defaults.Temperature, got.Temperature)
			assert.Equal(t, defaults.SafetySettings, request.SafetySettings)
			assert.Contains(t, requestPath, "/models/"+tt.model+":generateContent")
			if tt.name != "prompt example" {
				assert.Equal(t, "application/json", got.ResponseMIMEType)
				require.NotNil(t, got.ResponseJsonSchema)
			}

			afterJSON, err := json.Marshal(gc.gcfg)
			require.NoError(t, err)
			assert.JSONEq(t, string(defaultsJSON), string(afterJSON))
		})
	}
}
