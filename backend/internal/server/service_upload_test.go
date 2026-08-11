package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uploadedPlanRequest(planID, userID string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/uploads/"+planID, nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("plan_id", planID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	if userID != "" {
		ctx = context.WithValue(ctx, models.UserIdCtxKey, userID)
	}
	return req.WithContext(ctx)
}

func TestGetUploadedPlanHandlerRequiresAuthentication(t *testing.T) {
	called := false
	service := &RAGService{}
	lookup := func(context.Context, string, string) (*models.UploadedPlanResponse, error) {
		called = true
		return nil, nil
	}

	response := httptest.NewRecorder()
	service.getUploadedPlan(response, uploadedPlanRequest(uuid.NewString(), ""), lookup)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.False(t, called)
}

func TestGetUploadedPlanHandlerReturnsOwnedPlanWithoutUserID(t *testing.T) {
	planID := uuid.NewString()
	userID := uuid.NewString()
	service := &RAGService{}
	lookup := func(_ context.Context, gotPlanID, gotUserID string) (*models.UploadedPlanResponse, error) {
		assert.Equal(t, planID, gotPlanID)
		assert.Equal(t, userID, gotUserID)
		return &models.UploadedPlanResponse{
			PlanID:       planID,
			Title:        "Owner plan",
			Description:  "Private",
			Table:        models.Table{},
			AllowSharing: true,
		}, nil
	}

	response := httptest.NewRecorder()
	service.getUploadedPlan(response, uploadedPlanRequest(planID, userID), lookup)

	require.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `{"plan_id":"`+planID+`","created_at":"0001-01-01T00:00:00Z","title":"Owner plan","description":"Private","table":[],"allow_sharing":true}`, response.Body.String())
	assert.NotContains(t, response.Body.String(), "user_id")
}

func TestGetUploadedPlanHandlerHidesMissingAndUnauthorizedPlans(t *testing.T) {
	service := &RAGService{}
	lookup := func(context.Context, string, string) (*models.UploadedPlanResponse, error) {
		return nil, pgx.ErrNoRows
	}

	response := httptest.NewRecorder()
	service.getUploadedPlan(response, uploadedPlanRequest(uuid.NewString(), uuid.NewString()), lookup)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.Equal(t, "Plan not found\n", response.Body.String())
}

func TestGetUploadedPlanHandlerTreatsMalformedIDAsNotFound(t *testing.T) {
	called := false
	service := &RAGService{}
	lookup := func(context.Context, string, string) (*models.UploadedPlanResponse, error) {
		called = true
		return nil, nil
	}

	response := httptest.NewRecorder()
	service.getUploadedPlan(response, uploadedPlanRequest("not-a-uuid", uuid.NewString()), lookup)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.False(t, called)
}

func TestGetUploadedPlanHandlerHidesDatabaseErrors(t *testing.T) {
	service := &RAGService{}
	lookup := func(context.Context, string, string) (*models.UploadedPlanResponse, error) {
		return nil, errors.New("database password leaked")
	}

	response := httptest.NewRecorder()
	service.getUploadedPlan(response, uploadedPlanRequest(uuid.NewString(), uuid.NewString()), lookup)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Internal server error\n", response.Body.String())
	assert.False(t, strings.Contains(response.Body.String(), "database password"))
}
