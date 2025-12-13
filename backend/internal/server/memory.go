package server

import (
	"log/slog"
	"net/http"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
)

// DeleteMessageHandler handles the deletion of a single message from conversation history.
// The linked list structure will be repaired automatically.
// @Summary Delete a single message from conversation
// @Description Delete a specific message from the conversation history while maintaining the linked list integrity
// @Tags Memory
// @Accept json
// @Produce json
// @Param request body models.DeleteMessageRequest true "Request to delete a message"
// @Success 200 {string} string "Message deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /memory/message [delete]
func (rs *RAGService) DeleteMessageHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Deleting message...")

	// Get authenticated user ID
	userID, ok := req.Context().Value(models.UserIdCtxKey).(string)
	if !ok || userID == "" {
		logger.Error("User ID not found in context")
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	dmr := &models.DeleteMessageRequest{}
	err := models.GetRequestJSON(req, dmr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dmr.MessageID == "" {
		http.Error(w, "message_id is required", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "message_id", slog.StringValue(dmr.MessageID))

	err = rs.db.Memory.DeleteMessage(req.Context(), dmr.MessageID)
	if err != nil {
		logger.Error("Failed to delete message", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Message deleted successfully", "message_id", dmr.MessageID)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Message deleted successfully")); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// AddMessageHandler handles adding a new message to the conversation history.
// @Summary Add a message to conversation
// @Description Add a new message to the conversation history
// @Tags Memory
// @Accept json
// @Produce json
// @Param request body models.AddMessageRequest true "Request to add a message"
// @Success 200 {object} map[string]string "Message added successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /memory/message [post]
func (rs *RAGService) AddMessageHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Adding message...")

	// Get authenticated user ID
	userID, ok := req.Context().Value(models.UserIdCtxKey).(string)
	if !ok || userID == "" {
		logger.Error("User ID not found in context")
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	amr := &models.AddMessageRequest{}
	err := models.GetRequestJSON(req, amr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var prevMsgID *string
	if amr.PreviousMessageID != "" {
		prevMsgID = &amr.PreviousMessageID
	}

	msg, err := rs.db.Memory.AddMessage(req.Context(), amr.PlanID, userID, models.Role(amr.Role), amr.Content, prevMsgID, amr.PlanSnapshot.Plan())
	if err != nil {
		logger.Error("Failed to add message", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success and the new MessageID
	response := map[string]string{
		"message_id": msg.ID,
	}
	if err := models.WriteResponseJSON(w, http.StatusOK, response); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// DeleteMessagesAfterHandler handles the deletion of a message and all subsequent messages.
// This is useful for "branching" conversations by removing everything after a certain point.
// @Summary Delete a message and all subsequent messages
// @Description Delete a message and all messages that follow it in the conversation, allowing conversation branching
// @Tags Memory
// @Accept json
// @Produce json
// @Param request body models.DeleteMessagesAfterRequest true "Request to delete messages from a point"
// @Success 200 {string} string "Messages deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /memory/messages-after [delete]
func (rs *RAGService) DeleteMessagesAfterHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Deleting messages after specified message...")

	// Get authenticated user ID
	userID, ok := req.Context().Value(models.UserIdCtxKey).(string)
	if !ok || userID == "" {
		logger.Error("User ID not found in context")
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	dmar := &models.DeleteMessagesAfterRequest{}
	err := models.GetRequestJSON(req, dmar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dmar.MessageID == "" {
		http.Error(w, "message_id is required", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "message_id", slog.StringValue(dmar.MessageID))

	err = rs.db.Memory.DeleteMessagesAfter(req.Context(), dmar.MessageID)
	if err != nil {
		logger.Error("Failed to delete messages after", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Messages deleted successfully", "starting_from", dmar.MessageID)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Messages deleted successfully")); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// DeleteConversationHandler handles the deletion of an entire conversation.
// This removes all messages associated with a plan_id.
// @Summary Delete an entire conversation
// @Description Delete all messages in a conversation for a specific plan
// @Tags Memory
// @Accept json
// @Produce json
// @Param request body models.DeleteConversationRequest true "Request to delete a conversation"
// @Success 200 {string} string "Conversation deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /memory/conversation [delete]
func (rs *RAGService) DeleteConversationHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Deleting conversation...")

	// Get authenticated user ID
	userID, ok := req.Context().Value(models.UserIdCtxKey).(string)
	if !ok || userID == "" {
		logger.Error("User ID not found in context")
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	dcr := &models.DeleteConversationRequest{}
	err := models.GetRequestJSON(req, dcr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dcr.PlanID == "" {
		http.Error(w, "plan_id is required", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(dcr.PlanID))

	err = rs.db.Memory.DeleteConversation(req.Context(), dcr.PlanID)
	if err != nil {
		logger.Error("Failed to delete conversation", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Conversation deleted successfully", "plan_id", dcr.PlanID)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Conversation deleted successfully")); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}

// GetConversationHandler handles the retrieval of the conversation history for a plan.
// @Summary Get conversation history
// @Description Get the full conversation history for a specific plan
// @Tags Memory
// @Accept json
// @Produce json
// @Param plan_id query string true "Plan ID"
// @Success 200 {array} models.MessagePayload "Conversation history"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /memory/conversation [get]
func (rs *RAGService) GetConversationHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Getting conversation history...")

	// Get authenticated user ID
	userID, ok := req.Context().Value(models.UserIdCtxKey).(string)
	if !ok || userID == "" {
		logger.Error("User ID not found in context")
		http.Error(w, "Unauthorized: User ID missing", http.StatusUnauthorized)
		return
	}

	planID := req.URL.Query().Get("plan_id")
	if planID == "" {
		http.Error(w, "plan_id is required", http.StatusBadRequest)
		return
	}
	httplog.LogEntrySetField(req.Context(), "plan_id", slog.StringValue(planID))

	messages, err := rs.db.Memory.GetConversation(req.Context(), planID)
	if err != nil {
		logger.Error("Failed to get conversation", httplog.ErrAttr(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert messages to MessagePayloads
	var messagePayloads []models.MessagePayload
	for _, msg := range messages {
		msgPayload := models.MessagePayload{
			ID:                msg.ID,
			PlanID:            msg.PlanID,
			UserID:            msg.UserID,
			Role:              msg.Role,
			Content:           msg.Content,
			PreviousMessageID: msg.PreviousMessageID,
			NextMessageID:     msg.NextMessageID,
			PlanSnapshot:      nil,
			CreatedAt:         msg.CreatedAt,
		}
		if msg.PlanSnapshot != nil {
			msgPayload.PlanSnapshot = &models.RAGResponse{
				Title:       msg.PlanSnapshot.Title,
				Description: msg.PlanSnapshot.Description,
				PlanID:      msg.PlanSnapshot.PlanID,
				Table:       msg.PlanSnapshot.Table,
			}
		}
		messagePayloads = append(messagePayloads, msgPayload)

	}

	logger.Info("Conversation retrieved successfully", "plan_id", planID, "count", len(messages))
	if err := models.WriteResponseJSON(w, http.StatusOK, messagePayloads); err != nil {
		logger.Error("Failed to write response", httplog.ErrAttr(err))
	}
}
