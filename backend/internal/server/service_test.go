package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHandlerRejectsUnknownFields(t *testing.T) {
	service := &RAGService{}
	request := memoryHandlerRequest(http.MethodPost, "/chat", `{"plan_id":"00000000-0000-0000-0000-000000000001","message":"hello","forged":"value"}`, "user")
	response := httptest.NewRecorder()

	service.ChatHandler(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "invalid request body\n", response.Body.String())
}

func TestAddPlanToHistoryHidesRequestDecodeErrors(t *testing.T) {
	service := &RAGService{}
	request := memoryHandlerRequest(http.MethodPost, "/add-plan-to-history", `{"unknown":"value"}`, "user")
	response := httptest.NewRecorder()

	service.addPlanToHistory(
		response,
		request,
		func(context.Context, *models.Plan, string) error {
			t.Fatal("plan should not be persisted for an invalid request")
			return nil
		},
		func(context.Context, string, string) error { return nil },
		func(context.Context, string, string, models.Role, string, *string, *models.Plan) (*models.Message, error) {
			return nil, nil
		},
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Bad request\n", response.Body.String())
}

func TestAddPlanToHistoryRejectsInvalidPayload(t *testing.T) {
	service := &RAGService{}
	oversizedTitle := strings.Repeat("t", 201)
	request := memoryHandlerRequest(http.MethodPost, "/add-plan-to-history", `{"title":"`+oversizedTitle+`","description":"desc","table":[]}`, "user")
	response := httptest.NewRecorder()

	service.addPlanToHistory(
		response,
		request,
		func(context.Context, *models.Plan, string) error {
			t.Fatal("plan should not be persisted for an invalid request")
			return nil
		},
		func(context.Context, string, string) error { return nil },
		func(context.Context, string, string, models.Role, string, *string, *models.Plan) (*models.Message, error) {
			return nil, nil
		},
	)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Bad request\n", response.Body.String())
}

func TestAddPlanToHistoryCleansUpWhenInitialMessageFails(t *testing.T) {
	assertAddPlanToHistoryCleanup(t, func(_ int) (*models.Message, error) {
		return nil, errors.New("initial user message failed")
	})
}

func TestAddPlanToHistoryCleansUpWhenInitialAIMessageFails(t *testing.T) {
	assertAddPlanToHistoryCleanup(t, func(call int) (*models.Message, error) {
		if call == 1 {
			return &models.Message{ID: "user-message"}, nil
		}
		return nil, errors.New("initial AI message failed")
	})
}

func assertAddPlanToHistoryCleanup(t *testing.T, addMessage func(int) (*models.Message, error)) {
	t.Helper()
	const userID = "user"
	var deletedPlanID, deletedUserID string
	messageCalls := 0
	service := &RAGService{}
	request := memoryHandlerRequest(http.MethodPost, "/add-plan-to-history", `{"plan_id":"source","title":"Title","description":"Description","table":[],"initial_message":"hello"}`, userID)
	response := httptest.NewRecorder()

	service.addPlanToHistory(
		response,
		request,
		func(context.Context, *models.Plan, string) error { return nil },
		func(_ context.Context, planID, receivedUserID string) error {
			deletedPlanID, deletedUserID = planID, receivedUserID
			return nil
		},
		func(ctx context.Context, planID, receivedUserID string, role models.Role, content string, previousMessageID *string, snapshot *models.Plan) (*models.Message, error) {
			messageCalls++
			assert.NotEmpty(t, ctx)
			assert.NotEmpty(t, planID)
			assert.Equal(t, userID, receivedUserID)
			assert.NotEmpty(t, content)
			if messageCalls == 1 {
				assert.Equal(t, models.RoleUser, role)
				assert.Nil(t, previousMessageID)
				assert.Nil(t, snapshot)
			} else {
				assert.Equal(t, models.RoleAI, role)
				assert.NotNil(t, previousMessageID)
				assert.NotNil(t, snapshot)
			}
			return addMessage(messageCalls)
		},
	)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	assert.NotEmpty(t, deletedPlanID)
	assert.Equal(t, userID, deletedUserID)
	assert.GreaterOrEqual(t, messageCalls, 1)
}
