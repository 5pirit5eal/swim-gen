package server

import (
	"net/http"

	"github.com/5pirit5eal/swim-gen/internal/genai"
)

// rateLimitErrorMessage is returned to clients when the upstream AI provider
// responds with 429 RESOURCE_EXHAUSTED.
const rateLimitErrorMessage = "The AI service is currently overloaded. Please try again later."

// writeServiceError writes an HTTP error response for a failure from the
// RAG/LLM backend. Rate limiting (Vertex AI 429 RESOURCE_EXHAUSTED) is
// surfaced as HTTP 429 so the frontend can retry with backoff; any other
// error is reported as 500 with the fallback message.
func writeServiceError(w http.ResponseWriter, err error, fallbackMessage string) {
	if genai.IsRateLimitError(err) {
		http.Error(w, rateLimitErrorMessage, http.StatusTooManyRequests)
		return
	}
	http.Error(w, fallbackMessage, http.StatusInternalServerError)
}
