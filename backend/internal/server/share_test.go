package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sharePlanRequest(body, userID string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/share-plan", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if userID == "" {
		return req
	}
	return req.WithContext(context.WithValue(req.Context(), models.UserIdCtxKey, userID))
}

func TestSharePlanHandlerRequiresAuthentication(t *testing.T) {
	called := false
	service := &RAGService{}
	share := func(context.Context, string, string, models.SharingMethod) (string, error) {
		called = true
		return "hash", nil
	}

	response := httptest.NewRecorder()
	service.sharePlan(response, sharePlanRequest(`{"plan_id":"`+uuid.NewString()+`","method":"link"}`, ""), share)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.False(t, called)
}

func TestSharePlanHandlerValidatesRequest(t *testing.T) {
	service := &RAGService{}
	share := func(context.Context, string, string, models.SharingMethod) (string, error) {
		t.Fatal("share operation should not be called")
		return "", nil
	}

	tests := []struct {
		name string
		body string
	}{
		{name: "empty plan ID", body: `{"plan_id":"","method":"link"}`},
		{name: "malformed plan ID", body: `{"plan_id":"not-a-uuid","method":"link"}`},
		{name: "empty method", body: `{"plan_id":"` + uuid.NewString() + `","method":""}`},
		{name: "unsupported method", body: `{"plan_id":"` + uuid.NewString() + `","method":"email"}`},
		{name: "unknown field", body: `{"plan_id":"` + uuid.NewString() + `","method":"link","owner":"other"}`},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			service.sharePlan(response, sharePlanRequest(testCase.body, uuid.NewString()), share)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "Bad request\n", response.Body.String())
		})
	}
}

func TestSharePlanHandlerReturnsShareHash(t *testing.T) {
	planID := uuid.NewString()
	userID := uuid.NewString()
	var gotPlanID, gotUserID string
	var gotMethod models.SharingMethod
	service := &RAGService{}
	share := func(_ context.Context, receivedPlanID, receivedUserID string, method models.SharingMethod) (string, error) {
		gotPlanID, gotUserID, gotMethod = receivedPlanID, receivedUserID, method
		return "share-hash", nil
	}

	response := httptest.NewRecorder()
	service.sharePlan(response, sharePlanRequest(`{"plan_id":"`+planID+`","method":"link"}`, userID), share)

	require.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `{"url_hash":"share-hash"}`, response.Body.String())
	assert.Equal(t, planID, gotPlanID)
	assert.Equal(t, userID, gotUserID)
	assert.Equal(t, models.SharingMethodLink, gotMethod)
}

func TestSharePlanHandlerHidesUnauthorizedPlans(t *testing.T) {
	service := &RAGService{}
	share := func(context.Context, string, string, models.SharingMethod) (string, error) {
		return "", rag.ErrShareNotFound
	}

	response := httptest.NewRecorder()
	service.sharePlan(response, sharePlanRequest(`{"plan_id":"`+uuid.NewString()+`","method":"link"}`, uuid.NewString()), share)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "Plan not found\n", response.Body.String())
	assert.NotContains(t, response.Body.String(), "shareable plan not found")
}

func TestSharePlanHandlerHidesStoreErrors(t *testing.T) {
	service := &RAGService{}
	share := func(context.Context, string, string, models.SharingMethod) (string, error) {
		return "", errors.New("database password")
	}

	response := httptest.NewRecorder()
	service.sharePlan(response, sharePlanRequest(`{"plan_id":"`+uuid.NewString()+`","method":"link"}`, uuid.NewString()), share)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Internal server error\n", response.Body.String())
	assert.NotContains(t, response.Body.String(), "database password")
}
