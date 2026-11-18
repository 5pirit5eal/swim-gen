package server

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

func TestWriteServiceErrorMapsRateLimitTo429(t *testing.T) {
	err := fmt.Errorf("error generating plan: %w", genai.APIError{
		Code:    http.StatusTooManyRequests,
		Message: "Resource exhausted. Please try again later.",
		Status:  "RESOURCE_EXHAUSTED",
	})

	recorder := httptest.NewRecorder()
	writeServiceError(recorder, err, "Internal server error")

	assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
	assert.Equal(t, rateLimitErrorMessage+"\n", recorder.Body.String())
}

func TestWriteServiceErrorMapsOtherErrorsTo500(t *testing.T) {
	recorder := httptest.NewRecorder()
	writeServiceError(recorder, errors.New("database is down"), "Internal server error")

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, "Internal server error\n", recorder.Body.String())
}
