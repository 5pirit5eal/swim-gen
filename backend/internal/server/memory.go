package server

import (
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
