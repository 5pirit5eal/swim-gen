package genai

import (
	"errors"
	"net/http"

	"google.golang.org/genai"
)

// IsRateLimitError reports whether err is a Vertex AI API error with
// HTTP 429 / RESOURCE_EXHAUSTED, possibly wrapped by intermediate layers.
func IsRateLimitError(err error) bool {
	var apiErr genai.APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == http.StatusTooManyRequests || apiErr.Status == "RESOURCE_EXHAUSTED"
	}
	var apiErrPtr *genai.APIError
	if errors.As(err, &apiErrPtr) && apiErrPtr != nil {
		return apiErrPtr.Code == http.StatusTooManyRequests || apiErrPtr.Status == "RESOURCE_EXHAUSTED"
	}
	return false
}
