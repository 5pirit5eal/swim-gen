package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupabaseAuthMiddlewareReturnsGenericErrorForMalformedHeader(t *testing.T) {
	service := &RAGService{}
	nextCalled := false
	handler := service.SupabaseAuthMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		nextCalled = true
	}))
	req := httptest.NewRequest(http.MethodGet, "/uploads", nil)
	req.Header.Set("Authorization", "Basic sensitive-credentials")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Equal(t, "Unauthorized\n", response.Body.String())
	assert.NotContains(t, response.Body.String(), "Basic")
	assert.False(t, nextCalled)
}
