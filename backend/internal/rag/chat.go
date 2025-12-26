package rag

import (
	"context"
	"fmt"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
)

// ChatWithContext is the main stateless chat method for plan refinement through conversation.
// It retrieves conversation history from memory, builds context, calls the LLM, and stores the interaction.
func (db *RAGDB) ChatWithContext(
	ctx context.Context,
	planID, userID, userMessage string,
	lang models.Language,
	poolLength any,
) (*models.Plan, *models.Message, error) {
	logger := httplog.LogEntry(ctx)
	logger.Debug("Starting chat interaction", "plan_id", planID, "user_id", userID)

	// Require planID - application flow ensures plan exists before chat
	if planID == "" {
		logger.Error("Missing required plan_id for chat interaction")
		return nil, nil, fmt.Errorf("plan_id is required for chat interaction")
	}

	// 1. Retrieve conversation history (limited by config)
	conversation, err := db.Memory.GetConversation(ctx, planID)
	if err != nil {
		logger.Error("Failed to retrieve conversation history", httplog.ErrAttr(err))
		return nil, nil, fmt.Errorf("failed to retrieve conversation: %w", err)
	}

	// Apply history limit
	if len(conversation) > db.cfg.Chat.HistoryLimit {
		// Keep only the most recent N messages
		conversation = conversation[len(conversation)-db.cfg.Chat.HistoryLimit:]
		logger.Debug("Applied history limit", "total_messages", len(conversation), "limit", db.cfg.Chat.HistoryLimit)
	}

	// 2. Get current plan state
	var currentPlan *models.Plan
	plan, err := db.GetPlan(ctx, planID, SourceOptionPlan)
	if err == nil {
		currentPlan = plan.Plan()
		logger.Debug("Retrieved existing plan", "plan_id", planID, "title", currentPlan.Title)
	} else {
		// If plan doesn't exist, this is an invalid state or id
		logger.Error("Plan not found despite existing conversation", "plan_id", planID, httplog.ErrAttr(err))
		return nil, nil, fmt.Errorf("plan must exist for chat interaction: %w", err)
	}

	// 3. Build context
	conversationHistory := formatConversationHistory(conversation)
	contextDocs, err := db.buildChatContext(ctx, userMessage, currentPlan)
	if err != nil {
		logger.Warn("Failed to build chat context, continuing without reference docs", httplog.ErrAttr(err))
		contextDocs = []schema.Document{} // Continue without context if retrieval fails
	}

	// 4. Call GenAI ChatRefine
	db.Client.QueryMode() // Set embedder to query mode for any similarity searches
	chatResponse, err := db.Client.ChatRefine(
		ctx,
		conversationHistory,
		currentPlan,
		userMessage,
		string(lang),
		poolLength,
		contextDocs,
	)
	if err != nil {
		logger.Error("Failed to generate chat response", httplog.ErrAttr(err))
		return nil, nil, fmt.Errorf("failed to generate chat response: %w", err)
	}

	// 5. Store user message in memory
	var lastMessageID *string
	if len(conversation) > 0 {
		lastMessageID = &conversation[len(conversation)-1].ID
	}

	userMsg, err := db.Memory.AddMessage(
		ctx,
		planID,
		userID,
		models.RoleUser,
		userMessage,
		lastMessageID,
		nil, // No plan snapshot for user messages
	)
	if err != nil {
		logger.Error("Failed to store user message", httplog.ErrAttr(err))
		return nil, nil, fmt.Errorf("failed to store user message: %w", err)
	}

	// 6. Create plan from response and store AI message with plan snapshot
	var updatedPlan *models.Plan
	if chatResponse.Plan != nil {
		updatedPlan = &models.Plan{
			PlanID:      planID,
			Title:       chatResponse.Plan.Title,
			Description: chatResponse.Plan.Description,
			Table:       chatResponse.Plan.Table,
		}

		// Upsert plan to plans table
		_, err = db.UpsertPlan(ctx, *updatedPlan, userID)
		if err != nil {
			logger.Error("Failed to upsert plan", httplog.ErrAttr(err))
			return nil, nil, fmt.Errorf("failed to upsert plan: %w", err)
		}
	} else {
		// No plan update, use current plan if exists
		updatedPlan = currentPlan
	}

	// Store AI response with plan snapshot
	aiMsg, err := db.Memory.AddMessage(
		ctx,
		planID,
		userID,
		models.RoleAI,
		chatResponse.Response,
		&userMsg.ID,
		updatedPlan, // Store plan snapshot with AI response
	)
	if err != nil {
		logger.Error("Failed to store AI message", httplog.ErrAttr(err))
		return nil, nil, fmt.Errorf("failed to store AI message: %w", err)
	}

	logger.Debug("Chat interaction completed successfully", "plan_id", planID)
	return updatedPlan, aiMsg, nil
}

// buildChatContext retrieves similar plans from the vector store to provide reference context.
// This is called when additional context beyond the conversation history might be helpful.
func (db *RAGDB) buildChatContext(ctx context.Context, userQuery string, currentPlan *models.Plan) ([]schema.Document, error) {
	logger := httplog.LogEntry(ctx)

	// Check if RAG context is enabled in config
	if !db.cfg.Chat.UseRAGContext {
		logger.Debug("RAG context disabled in config")
		return []schema.Document{}, nil
	}

	// Retrieve context if no plan exists yet or if RAG context is enabled
	if currentPlan == nil || db.cfg.Chat.UseRAGContext {
		// Perform similarity search for relevant plans
		docs, err := db.PlanStore.SimilaritySearch(ctx, userQuery, 3)
		if err != nil {
			logger.Warn("Similarity search failed", httplog.ErrAttr(err))
			return []schema.Document{}, nil // Return empty, don't fail the whole request
		}
		logger.Debug("Retrieved reference plans", "count", len(docs))
		return docs, nil
	}

	// No additional context needed
	return []schema.Document{}, nil
}

// formatConversationHistory converts message history into a text format suitable for LLM prompts.
// It includes plan snapshots when the trainer created or modified plans.
func formatConversationHistory(messages []models.Message) string {
	if len(messages) == 0 {
		return "Noch keine Nachrichten im Gespräch."
	}

	var builder strings.Builder
	for _, msg := range messages {
		if msg.Role == models.RoleUser {
			builder.WriteString(fmt.Sprintf("Schwimmer: %s\\n", msg.Content))
		} else {
			// AI trainer response
			builder.WriteString(fmt.Sprintf("Trainer: %s\\n", msg.Content))

			// Include plan snapshot if the trainer created/modified a plan
			if msg.PlanSnapshot != nil {
				builder.WriteString("  [Plan erstellt/geändert]:\n")
				builder.WriteString(fmt.Sprintf("  %s\n", msg.PlanSnapshot.String()))
			}
		}
	}
	return builder.String()
}
