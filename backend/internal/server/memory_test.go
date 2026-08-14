package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type memoryHandlerFake struct {
	add                func(context.Context, string, string, models.Role, string, *string, *models.Plan) (*models.Message, error)
	get                func(context.Context, string, string) ([]models.Message, error)
	deleteConversation func(context.Context, string, string) error
	delete             func(context.Context, string, string) error
	deleteAfter        func(context.Context, string, string) error
}

func (f *memoryHandlerFake) AddMessage(ctx context.Context, planID, userID string, role models.Role, content string, previousMessageID *string, planSnapshot *models.Plan) (*models.Message, error) {
	return f.add(ctx, planID, userID, role, content, previousMessageID, planSnapshot)
}

func (f *memoryHandlerFake) GetConversation(ctx context.Context, planID, userID string) ([]models.Message, error) {
	return f.get(ctx, planID, userID)
}

func (*memoryHandlerFake) GetLastMessage(context.Context, pgxscan.Querier, string, string) (*models.Message, error) {
	return nil, nil
}

func (f *memoryHandlerFake) DeleteConversation(ctx context.Context, planID, userID string) error {
	return f.deleteConversation(ctx, planID, userID)
}

func (f *memoryHandlerFake) DeleteMessage(ctx context.Context, messageID, userID string) error {
	return f.delete(ctx, messageID, userID)
}

func (f *memoryHandlerFake) DeleteMessagesAfter(ctx context.Context, messageID, userID string) error {
	return f.deleteAfter(ctx, messageID, userID)
}

func memoryHandlerRequest(method, target, body, userID string) *http.Request {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if userID == "" {
		return request
	}
	return request.WithContext(context.WithValue(request.Context(), models.UserIdCtxKey, userID))
}

func TestMemoryHandlersRequireAuthentication(t *testing.T) {
	called := false
	fake := &memoryHandlerFake{
		delete: func(context.Context, string, string) error {
			called = true
			return nil
		},
	}
	service := &RAGService{db: &rag.RAGDB{Memory: fake}}

	response := httptest.NewRecorder()
	service.DeleteMessageHandler(response, memoryHandlerRequest(http.MethodDelete, "/memory/message", `{"message_id":"message"}`, ""))

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.False(t, called)
}

func TestAddMessageHandlerPassesAuthenticatedOwnerAndRejectsForgedAI(t *testing.T) {
	const userID = "owner"
	var gotUserID string
	var gotRole models.Role
	fake := &memoryHandlerFake{
		add: func(_ context.Context, _, receivedUserID string, role models.Role, _ string, _ *string, snapshot *models.Plan) (*models.Message, error) {
			gotUserID = receivedUserID
			gotRole = role
			assert.Nil(t, snapshot)
			return &models.Message{ID: "message"}, nil
		},
	}
	service := &RAGService{db: &rag.RAGDB{Memory: fake}}

	response := httptest.NewRecorder()
	service.AddMessageHandler(response, memoryHandlerRequest(http.MethodPost, "/memory/message", `{"plan_id":"plan","content":"hello"}`, userID))

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, userID, gotUserID)
	assert.Equal(t, models.RoleUser, gotRole)

	response = httptest.NewRecorder()
	service.AddMessageHandler(response, memoryHandlerRequest(http.MethodPost, "/memory/message", `{"plan_id":"plan","role":"ai","content":"forged"}`, userID))

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestAddMessageHandlerHidesStoreErrors(t *testing.T) {
	fake := &memoryHandlerFake{
		add: func(context.Context, string, string, models.Role, string, *string, *models.Plan) (*models.Message, error) {
			return nil, errors.New("database password")
		},
	}
	service := &RAGService{db: &rag.RAGDB{Memory: fake}}

	response := httptest.NewRecorder()
	service.AddMessageHandler(response, memoryHandlerRequest(http.MethodPost, "/memory/message", `{"plan_id":"plan","content":"hello"}`, "owner"))

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Internal server error\n", response.Body.String())
	assert.NotContains(t, response.Body.String(), "database password")
}

func TestMemoryHandlersPassAuthenticatedUserToReadsAndDeletes(t *testing.T) {
	const userID = "owner"
	var gotPlanID, gotMessageID, gotUserID string
	fake := &memoryHandlerFake{
		get: func(_ context.Context, planID, receivedUserID string) ([]models.Message, error) {
			gotPlanID, gotUserID = planID, receivedUserID
			return []models.Message{{
				ID:        "message",
				PlanID:    planID,
				UserID:    receivedUserID,
				Role:      models.RoleUser,
				Content:   "hello",
				CreatedAt: time.Unix(0, 0).UTC(),
			}}, nil
		},
		delete: func(_ context.Context, messageID, receivedUserID string) error {
			gotMessageID, gotUserID = messageID, receivedUserID
			return nil
		},
		deleteAfter: func(context.Context, string, string) error { return nil },
		deleteConversation: func(context.Context, string, string) error {
			return nil
		},
	}
	service := &RAGService{db: &rag.RAGDB{Memory: fake}}

	response := httptest.NewRecorder()
	service.GetConversationHandler(response, memoryHandlerRequest(http.MethodGet, "/memory/conversation?plan_id=plan", "", userID))

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "plan", gotPlanID)
	assert.Equal(t, userID, gotUserID)
	assert.NotContains(t, response.Body.String(), "user_id")

	response = httptest.NewRecorder()
	service.DeleteMessageHandler(response, memoryHandlerRequest(http.MethodDelete, "/memory/message", `{"message_id":"message"}`, userID))

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "message", gotMessageID)
	assert.Equal(t, userID, gotUserID)

}

func TestMemoryHandlersHideUnauthorizedResources(t *testing.T) {
	fake := &memoryHandlerFake{
		get: func(context.Context, string, string) ([]models.Message, error) {
			return nil, rag.ErrMemoryNotFound
		},
		delete: func(context.Context, string, string) error {
			return rag.ErrMemoryNotFound
		},
		deleteAfter: func(context.Context, string, string) error {
			return rag.ErrMemoryNotFound
		},
		deleteConversation: func(context.Context, string, string) error {
			return rag.ErrMemoryNotFound
		},
	}
	service := &RAGService{db: &rag.RAGDB{Memory: fake}}

	cases := []struct {
		name    string
		method  string
		target  string
		body    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"read", http.MethodGet, "/memory/conversation?plan_id=plan", "", service.GetConversationHandler},
		{"single delete", http.MethodDelete, "/memory/message", `{"message_id":"message"}`, service.DeleteMessageHandler},
		{"chain delete", http.MethodDelete, "/memory/messages-after", `{"message_id":"message"}`, service.DeleteMessagesAfterHandler},
		{"conversation delete", http.MethodDelete, "/memory/conversation", `{"plan_id":"plan"}`, service.DeleteConversationHandler},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			testCase.handler(response, memoryHandlerRequest(testCase.method, testCase.target, testCase.body, "other-user"))
			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.NotContains(t, response.Body.String(), "memory resource not found")
		})
	}
}
