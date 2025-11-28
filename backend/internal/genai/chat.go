package genai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
	"google.golang.org/genai"
)

// ChatRefine generates or refines a training plan based on conversation context.
// It uses the conversation history, current plan state, and user's latest message
// to create or update a plan while maintaining conversational context.
func (gc *GoogleGenAIClient) ChatRefine(
	ctx context.Context,
	conversationHistory string,
	currentPlan *models.Plan,
	userMessage string,
	lang string,
	poolLength any,
	contextDocs []schema.Document,
) (*models.ChatResponse, error) {
	logger := httplog.LogEntry(ctx)

	// Get the ChatResponse JSON schema
	chatSchema, err := models.ChatResponseSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get ChatResponse schema: %w", err)
	}

	// Format reference plans for context
	var contextStr string
	if len(contextDocs) > 0 {
		for i, doc := range contextDocs {
			contextStr += fmt.Sprintf("Plan %d:\n%s\n\n", i+1, doc.PageContent)
		}
	} else {
		contextStr = "Keine Referenzpläne verfügbar."
	}

	// Format current plan
	var currentPlanStr string
	if currentPlan != nil {
		currentPlanStr = currentPlan.String()
	} else {
		currentPlanStr = "Noch kein Plan erstellt."
	}

	// Create the chat refinement query
	query := fmt.Sprintf(
		chatRefineTemplateStr,
		poolLength,
		lang,
		conversationHistory,
		currentPlanStr,
		contextStr,
		userMessage,
	)

	// Configure generation with structured output
	genCfg := *gc.gcfg
	genCfg.ResponseMIMEType = "application/json"
	genCfg.ResponseJsonSchema = chatSchema

	// Call the LLM
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	if err != nil {
		logger.Error("Error when generating chat response with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error when generating chat response: %w", err)
	}

	// Parse the response
	var chatResponse models.ChatResponse
	err = json.Unmarshal([]byte(answer.Text()), &chatResponse)
	if err != nil {
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err), "raw_response", answer.Text())
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}

	// If plan was generated, update sums
	if chatResponse.Plan != nil && len(chatResponse.Plan.Table) > 0 {
		// Ensure total row exists
		if !containsTotal(chatResponse.Plan.Table) {
			chatResponse.Plan.Table.AddSum()
		}
		// Recalculate sums to ensure correctness
		chatResponse.Plan.Table.UpdateSum()
	}

	logger.Debug("Chat response generated successfully")
	return &chatResponse, nil
}

// containsTotal checks if the table contains a total row
func containsTotal(table models.Table) bool {
	if len(table) == 0 {
		return false
	}
	lastRow := table[len(table)-1]
	return lastRow.Content == "Gesamt" || lastRow.Content == "Total"
}
