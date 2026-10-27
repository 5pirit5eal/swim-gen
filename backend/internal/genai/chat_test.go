package genai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tmc/langchaingo/schema"
	"google.golang.org/genai"
)

func TestChatRefinePreservesStructuredPlansAndGrounding(t *testing.T) {
	currentPlan := &models.Plan{
		Title: "Technique session", Description: "Freestyle with fins",
		Table: models.Table{{Amount: 4, Content: "Technique set", SubRows: []models.Row{
			{Amount: 1, Distance: 50, Content: "[Catch-up](/drill/catch-up)", Intensity: "TÜ", Equipment: []models.EquipmentType{models.EquipmentFins}},
			{Amount: 1, Distance: 50, Content: "Freestyle", Intensity: "02:00 / 100m", Break: "Rest until 2:15"},
		}}},
	}
	drills := []schema.Document{
		{PageContent: "Catch-up drill instructions", Metadata: map[string]any{"slug": "Catch-up", "img_name": "catch-up.webp"}},
		{PageContent: "Sculling drill instructions", Metadata: map[string]any{"slug": "Sculling", "img_name": "sculling.png"}},
	}

	for _, noChange := range []bool{false, true} {
		name := "refinement"
		if noChange {
			name = "question without plan changes"
		}
		t.Run(name, func(t *testing.T) {
			response := models.ChatResponse{Response: "Trainer response"}
			if !noChange {
				response.Plan = &models.GeneratedPlan{Title: currentPlan.Title, Description: currentPlan.Description, Table: currentPlan.Table}
			}
			responseJSON, err := json.Marshal(response)
			require.NoError(t, err)

			var request struct {
				Contents         []*genai.Content `json:"contents"`
				GenerationConfig struct {
					ResponseMIMEType   string         `json:"responseMimeType"`
					ResponseJSONSchema map[string]any `json:"responseJsonSchema"`
				} `json:"generationConfig"`
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.NoError(t, json.NewDecoder(r.Body).Decode(&request))
				w.Header().Set("Content-Type", "application/json")
				assert.NoError(t, json.NewEncoder(w).Encode(&genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{{Content: genai.NewContentFromText(string(responseJSON), genai.RoleModel)}},
				}))
			}))
			defer server.Close()
			client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
				Backend: genai.BackendVertexAI, Project: "test", Location: "test",
				HTTPClient: server.Client(), HTTPOptions: genai.HTTPOptions{BaseURL: server.URL},
			})
			require.NoError(t, err)
			gc := &GoogleGenAIClient{gc: client, gcfg: &genai.GenerateContentConfig{}}
			gc.cfg.Model = "test-model"

			result, err := gc.ChatRefine(context.Background(), "Previous request", currentPlan, "Latest request", "en", 25, nil, drills)
			require.NoError(t, err)
			assert.Equal(t, "application/json", request.GenerationConfig.ResponseMIMEType)
			expectedSchema, err := models.ChatResponseSchema()
			require.NoError(t, err)
			assert.Equal(t, expectedSchema, request.GenerationConfig.ResponseJSONSchema)
			require.Len(t, request.Contents, 1)
			prompt := request.Contents[0].Parts[0].Text
			assert.Contains(t, prompt, planFormatInstructions)
			assert.Contains(t, ragTemplateStr, planFormatInstructions)
			assert.Contains(t, prompt, "Catch-up drill instructions\n\nReferenz: [Catch-up](/drill/catch-up)")
			assert.Contains(t, prompt, "Sculling drill instructions\n\nReferenz: [Sculling](/drill/sculling)")
			assert.NotContains(t, prompt, "%!")

			_, planSection, found := strings.Cut(prompt, "AKTUELLER PLAN (JSON):\n")
			require.True(t, found)
			var suppliedPlan models.GeneratedPlan
			require.NoError(t, json.NewDecoder(strings.NewReader(planSection)).Decode(&suppliedPlan))
			assert.Equal(t, currentPlan.Table, suppliedPlan.Table)

			assert.Equal(t, response.Response, result.Response)
			if noChange {
				assert.Nil(t, result.Plan)
				return
			}
			require.NotNil(t, result.Plan)
			require.Len(t, result.Plan.Table, 2)
			set := result.Plan.Table[0]
			assert.Equal(t, 100, set.Distance)
			assert.Equal(t, 400, set.Sum)
			assert.Equal(t, 400, result.Plan.Table[1].Sum)
			assert.Equal(t, currentPlan.Table[0].SubRows[0].Content, set.SubRows[0].Content)
			assert.Equal(t, currentPlan.Table[0].SubRows[0].Equipment, set.SubRows[0].Equipment)
			assert.Equal(t, currentPlan.Table[0].SubRows[1].Intensity, set.SubRows[1].Intensity)
		})
	}
}
