package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/invopop/jsonschema"
)

type UserProfile struct {
	UserID             string    `db:"user_id"`
	UpdatedAt          time.Time `db:"updated_at"`
	Username           string    `db:"username"`
	Experience         *string   `db:"experience,omitempty"`
	PreferredLanguage  *string   `db:"preferred_language,omitempty"`
	PreferredStrokes   []string  `db:"preferred_strokes"`
	Categories         []string  `db:"categories"`
	OverallGenerations int       `db:"overall_generations"`
	MonthlyGenerations int       `db:"monthly_generations"`
	Exports            int       `db:"exports"`
}

type Feedback struct {
	UserID           string    `db:"user_id"`
	PlanID           string    `db:"plan_id"`
	Rating           int       `db:"rating"`
	WasSwam          bool      `db:"was_swam"`
	DifficultyRating int       `db:"difficulty_rating"`
	Comment          string    `db:"comment"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type ChoiceResult struct {
	Idx         int    `json:"index" example:"1"`
	Description string `json:"description" example:"Selected plan based on your requirements"`
}

// ChatResponse represents the structured response from a chat interaction
// @Description Response containing both a training plan and conversational response from the chat system
type ChatResponse struct {
	Plan     *GeneratedPlan `json:"plan" jsonschema_description:"The training plan created or refined based on the conversation"`
	Response string         `json:"response" jsonschema_description:"Conversational response explaining the plan or answering the user's question"`
}

// ChatResponseSchema generates the JSON schema for ChatResponse
func ChatResponseSchema() (map[string]any, error) {
	schema := jsonschema.Reflect(&ChatResponse{})

	jsonSchema, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON schema: %w", err)
	}
	var result map[string]any
	return result, json.Unmarshal(jsonSchema, &result)
}
