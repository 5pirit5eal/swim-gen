package server

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/rag"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func deletePlanRequest(planID, userID string) *http.Request {
	req := httptest.NewRequest(http.MethodDelete, "/plan/"+planID, nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("plan_id", planID)
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	if userID != "" {
		ctx = context.WithValue(ctx, models.UserIdCtxKey, userID)
	}
	return req.WithContext(ctx)
}

func TestDeletePlanHandlerRequiresAuthentication(t *testing.T) {
	called := false
	service := &RAGService{}
	deleteFn := func(context.Context, string, string) error {
		called = true
		return nil
	}

	response := httptest.NewRecorder()
	service.deletePlan(response, deletePlanRequest(uuid.NewString(), ""), deleteFn)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.False(t, called)
}

func TestDeletePlanHandlerRejectsMalformedPlanID(t *testing.T) {
	called := false
	service := &RAGService{}
	deleteFn := func(context.Context, string, string) error {
		called = true
		return nil
	}

	response := httptest.NewRecorder()
	service.deletePlan(response, deletePlanRequest("invalid-uuid", uuid.NewString()), deleteFn)

	assert.Equal(t, http.StatusNotFound, response.Code)
	assert.False(t, called)
}

func TestDeletePlanHandlerHidesUnauthorizedOrMissingPlans(t *testing.T) {
	service := &RAGService{}
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "err plan not found",
			err:  rag.ErrPlanNotFound,
		},
		{
			name: "err no rows",
			err:  pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleteFn := func(context.Context, string, string) error {
				return tt.err
			}

			response := httptest.NewRecorder()
			service.deletePlan(response, deletePlanRequest(uuid.NewString(), uuid.NewString()), deleteFn)

			assert.Equal(t, http.StatusNotFound, response.Code)
			assert.Equal(t, "Plan not found\n", response.Body.String())
		})
	}
}

func TestDeletePlanHandlerHidesDatabaseErrors(t *testing.T) {
	service := &RAGService{}
	deleteFn := func(context.Context, string, string) error {
		return errors.New("connection failed: sensitive db connection string")
	}

	response := httptest.NewRecorder()
	service.deletePlan(response, deletePlanRequest(uuid.NewString(), uuid.NewString()), deleteFn)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "Internal server error\n", response.Body.String())
	assert.NotContains(t, response.Body.String(), "sensitive")
}

func TestDeletePlanHandlerDeletesOwnedPlanSuccessfully(t *testing.T) {
	planID := uuid.NewString()
	userID := uuid.NewString()
	called := false

	service := &RAGService{}
	deleteFn := func(_ context.Context, gotPlanID, gotUserID string) error {
		called = true
		assert.Equal(t, planID, gotPlanID)
		assert.Equal(t, userID, gotUserID)
		return nil
	}

	response := httptest.NewRecorder()
	service.deletePlan(response, deletePlanRequest(planID, userID), deleteFn)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.True(t, called)
	assert.Equal(t, "Plan deleted successfully", response.Body.String())
}
