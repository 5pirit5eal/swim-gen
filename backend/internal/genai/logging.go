package genai

import (
	"log/slog"
	"time"
)

// logInference records only the model call, excluding prompt preparation and response parsing.
func logInference(logger *slog.Logger, operation, model string) func(error) {
	logger = logger.With("operation", operation, "model", model)
	logger.Info("Starting model inference")
	startedAt := time.Now()
	return func(err error) {
		logger.Info("Model inference finished", "duration_ms", time.Since(startedAt).Milliseconds(), "success", err == nil)
	}
}
