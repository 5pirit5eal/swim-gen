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
	contextDocs, drillDocs []schema.Document,
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
		planJSON, err := json.Marshal(models.GeneratedPlan{
			Title: currentPlan.Title, Description: currentPlan.Description, Table: currentPlan.Table,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to marshal current plan: %w", err)
		}
		currentPlanStr = string(planJSON)
	} else {
		currentPlanStr = "null"
	}

	// Create the chat refinement query
	query := fmt.Sprintf(
		chatRefineTemplateStr,
		poolLength,
		lang,
		conversationHistory,
		currentPlanStr,
		contextStr,
		formatDrillDocuments(drillDocs),
		userMessage,
	)

	// Configure generation with structured output
	genCfg := *gc.gcfg
	genCfg.ThinkingConfig = &genai.ThinkingConfig{
		IncludeThoughts: false,
		ThinkingLevel:   genai.ThinkingLevelLow,
	}
	genCfg.ResponseMIMEType = "application/json"
	genCfg.ResponseJsonSchema = chatSchema

	// Call the LLM
	finishInference := logInference(logger, "chat_refine", gc.cfg.Model)
	answer, err := gc.gc.Models.GenerateContent(ctx, gc.cfg.Model, genai.Text(query), &genCfg)
	finishInference(err)
	if err != nil {
		logger.Error("Error when generating chat response with LLM", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error when generating chat response: %w", err)
	}

	// Parse the response
	var chatResponse models.ChatResponse
	err = json.Unmarshal([]byte(answer.Text()), &chatResponse)
	if err != nil {
		logger.Debug("LLM response could not be parsed", "raw_response", answer.Text())
		logger.Error("Error parsing LLM response", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error parsing LLM response: %w", err)
	}

	// If plan was generated, update sums
	if chatResponse.Plan != nil && len(chatResponse.Plan.Table) > 0 {
		chatResponse.Plan.Table.FlattenSingleParentRow()
		// Ensure total row exists
		chatResponse.Plan.Table.AddSum()
		// Recalculate sums to ensure correctness
		chatResponse.Plan.Table.UpdateSum()
	}

	logger.Debug("Chat response generated successfully")
	return &chatResponse, nil
}
