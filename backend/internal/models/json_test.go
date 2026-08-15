package models_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sampleJSONPayload struct {
	Name string `json:"name"`
}

func TestGetRequestJSON(t *testing.T) {
	t.Run("Valid JSON within limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(`{"name":"swim"}`)))
		req.Header.Set("Content-Type", "application/json")

		var payload sampleJSONPayload
		err := models.GetRequestJSON(req, &payload)
		require.NoError(t, err)
		assert.Equal(t, "swim", payload.Name)
	})

	t.Run("Rejects unsupported content type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte(`{"name":"swim"}`)))
		req.Header.Set("Content-Type", "text/plain")

		var payload sampleJSONPayload
		err := models.GetRequestJSON(req, &payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported content type")
	})

	t.Run("Rejects oversized JSON payload", func(t *testing.T) {
		// Create a JSON payload larger than 2MB
		oversized := `{"name":"` + strings.Repeat("a", models.MaxJSONBodyBytes+100) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(oversized))
		req.Header.Set("Content-Type", "application/json")

		var payload sampleJSONPayload
		err := models.GetRequestJSON(req, &payload)
		assert.Error(t, err)
	})
}
