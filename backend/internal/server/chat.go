package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
)

// ChatHandler handles chat-based training plan creation and refinement.
// It processes conversational interactions to iteratively build and improve training plans.
// @Summary Chat with AI to create or refine training plans
// @Description Have a conversation with the AI trainer to create, modify, or get information about training plans
// @Tags Chat
// @Accept json
// @Produce json
// @Param request body models.ChatRequest true "Chat request with message and optional plan ID"
// @Success 200 {object} models.ChatResponsePayload "Successful chat response with updated plan"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /chat [post]
func (rs *RAGService) ChatHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())

	// Get authenticated user ID
	userID, ok := req.Context().Value(models.UserIdCtxKey).(string)
	if !ok {
		logger.Error("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger.Info("Processing chat request...", "user_id", userID)
	// Parse request
	var chatReq models.ChatRequest
	if err := json.NewDecoder(req.Body).Decode(&chatReq); err != nil {
		logger.Error("Failed to decode request body", httplog.ErrAttr(err))
		http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}
	if chatReq.PlanID != "" {
		httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(chatReq.PlanID))
	}

	// Set defaults
	if chatReq.Language == "" {
		chatReq.Language = models.LanguageDE
	}
	if chatReq.PoolLength == nil {
		chatReq.PoolLength = 25
	}

	logger.Debug("Chat request parsed",
		"plan_id", chatReq.PlanID,
		"message", chatReq.Message,
		"language", chatReq.Language,
		"pool_length", chatReq.PoolLength,
	)

	// Call ChatWithContext
	updatedPlan, aiMessage, err := rs.db.ChatWithContext(
		req.Context(),
		chatReq.PlanID,
		userID,
		chatReq.Message,
		chatReq.Language,
		chatReq.PoolLength,
	)
	if err != nil {
		logger.Error("Failed to process chat interaction", httplog.ErrAttr(err))
		http.Error(w, fmt.Sprintf("failed to process chat: %v", err), http.StatusInternalServerError)
		return
	}

	// Build response
	response := models.ChatResponsePayload{
		PlanID:   chatReq.PlanID,
		Response: aiMessage.Content,
	}

	// Include plan details if plan was created/updated
	if updatedPlan != nil {
		response.Title = updatedPlan.Title
		response.Description = updatedPlan.Description
		response.Table = updatedPlan.Table
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("Failed to encode response", httplog.ErrAttr(err))
	}

	logger.Info("Chat request completed successfully", "user_id", userID, "plan_id", response.PlanID)
}
