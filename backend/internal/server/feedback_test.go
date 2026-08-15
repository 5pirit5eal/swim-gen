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

func feedbackRequest(body, userID string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if userID == "" {
		return req
	}
	return req.WithContext(context.WithValue(req.Context(), models.UserIdCtxKey, userID))
}

func validFeedbackBody(planID string) string {
	return `{"plan_id":"` + planID + `","rating":5,"was_swam":true,"difficulty_rating":7,"comment":"Useful"}`
}

func TestFeedbackHandlerRequiresAuthentication(t *testing.T) {
	called := false
	service := &RAGService{}
	submit := func(context.Context, *models.Feedback) error {
		called = true
		return nil
	}

	response := httptest.NewRecorder()
	service.submitFeedback(response, feedbackRequest(validFeedbackBody(uuid.NewString()), ""), submit)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.False(t, called)
}

func TestFeedbackHandlerValidatesRequest(t *testing.T) {
	service := &RAGService{}
	submit := func(context.Context, *models.Feedback) error {
		t.Fatal("feedback operation should not be called")
		return nil
	}

	planID := uuid.NewString()
	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"plan_id":`},
		{name: "unknown field", body: validFeedbackBody(planID)[:len(validFeedbackBody(planID))-1] + `,"user_id":"other"}`},
		{name: "malformed plan ID", body: validFeedbackBody("not-a-uuid")},
		{name: "missing plan ID", body: validFeedbackBody("")},
		{name: "rating below range", body: `{"plan_id":"` + planID + `","rating":0,"difficulty_rating":7}`},
		{name: "rating above range", body: `{"plan_id":"` + planID + `","rating":6,"difficulty_rating":7}`},
		{name: "difficulty below range", body: `{"plan_id":"` + planID + `","rating":5,"difficulty_rating":0}`},
		{name: "difficulty above range", body: `{"plan_id":"` + planID + `","rating":5,"difficulty_rating":11}`},
		{name: "comment too long", body: `{"plan_id":"` + planID + `","rating":5,"difficulty_rating":5,"comment":"` + strings.Repeat("c", 1001) + `"}`},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			service.submitFeedback(response, feedbackRequest(testCase.body, uuid.NewString()), submit)

			assert.Equal(t, http.StatusBadRequest, response.Code)
			assert.Equal(t, "Bad request\n", response.Body.String())
		})
	}
}

func TestFeedbackHandlerSubmitsForAuthenticatedUser(t *testing.T) {
	planID := uuid.NewString()
	userID := uuid.NewString()
	var received *models.Feedback
	service := &RAGService{}
	submit := func(_ context.Context, feedback *models.Feedback) error {
		received = feedback
		return nil
	}

	response := httptest.NewRecorder()
	service.submitFeedback(response, feedbackRequest(validFeedbackBody(planID), userID), submit)

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "Feedback submitted successfully", response.Body.String())
	require.NotNil(t, received)
	assert.Equal(t, userID, received.UserID)
	assert.Equal(t, planID, received.PlanID)
	assert.Equal(t, 5, received.Rating)
	assert.True(t, received.WasSwam)
	require.NotNil(t, received.DifficultyRating)
	assert.Equal(t, 7, *received.DifficultyRating)
	assert.Equal(t, "Useful", received.Comment)
}

func TestFeedbackHandlerHidesUnauthorizedPlans(t *testing.T) {
	service := &RAGService{}
	submit := func(context.Context, *models.Feedback) error {
		return rag.ErrFeedbackPlanNotFound
	}

	response := httptest.NewRecorder()
	service.submitFeedback(response, feedbackRequest(validFeedbackBody(uuid.NewString()), uuid.NewString()), submit)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "Plan not found\n", response.Body.String())
}

func TestFeedbackHandlerHidesStoreErrors(t *testing.T) {
	service := &RAGService{}
	submit := func(context.Context, *models.Feedback) error {
		return errors.New("database password leaked")
	}

	response := httptest.NewRecorder()
	service.submitFeedback(response, feedbackRequest(validFeedbackBody(uuid.NewString()), uuid.NewString()), submit)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Internal server error\n", response.Body.String())
	assert.NotContains(t, response.Body.String(), "database password")
}
