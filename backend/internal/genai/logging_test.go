package genai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogInference(t *testing.T) {
	for _, inferenceErr := range []error{nil, errors.New("model unavailable"), context.DeadlineExceeded} {
		name := "success"
		if inferenceErr != nil {
			name = inferenceErr.Error()
		}
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&output, nil)).With("request_id", "test-request")
			finishInference := logInference(logger, "generate_plan", "test-model")

			decoder := json.NewDecoder(&output)
			var start map[string]any
			require.NoError(t, decoder.Decode(&start))
			assert.Equal(t, "Starting model inference", start["msg"])
			assert.NotContains(t, start, "duration_ms")
			finishInference(inferenceErr)
			var finish map[string]any
			require.NoError(t, decoder.Decode(&finish))
			assert.Equal(t, "Model inference finished", finish["msg"])
			assert.Equal(t, inferenceErr == nil, finish["success"])
			require.Contains(t, finish, "duration_ms")
			assert.GreaterOrEqual(t, finish["duration_ms"].(float64), float64(0))
			for _, entry := range []map[string]any{start, finish} {
				assert.Equal(t, "INFO", entry["level"])
				assert.Equal(t, "generate_plan", entry["operation"])
				assert.Equal(t, "test-model", entry["model"])
				assert.Equal(t, "test-request", entry["request_id"])
				assert.NotContains(t, entry, "prompt")
				assert.NotContains(t, entry, "response")
			}
		})
	}
}
