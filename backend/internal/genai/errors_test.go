package genai

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genai"
)

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "unrelated error",
			err:  errors.New("some other failure"),
			want: false,
		},
		{
			name: "APIError with code 429",
			err: genai.APIError{
				Code:    http.StatusTooManyRequests,
				Message: "Resource exhausted. Please try again later.",
				Status:  "RESOURCE_EXHAUSTED",
			},
			want: true,
		},
		{
			name: "APIError with RESOURCE_EXHAUSTED status only",
			err:  genai.APIError{Status: "RESOURCE_EXHAUSTED"},
			want: true,
		},
		{
			name: "wrapped APIError with code 429",
			err: fmt.Errorf("error generating plan: %w",
				fmt.Errorf("Models.GenerateContent: %w", genai.APIError{
					Code:    429,
					Message: "Resource exhausted. Please try again later.",
					Status:  "RESOURCE_EXHAUSTED",
				})),
			want: true,
		},
		{
			name: "pointer APIError with code 429",
			err:  fmt.Errorf("failed: %w", &genai.APIError{Code: 429, Status: "RESOURCE_EXHAUSTED"}),
			want: true,
		},
		{
			name: "APIError with different code",
			err:  genai.APIError{Code: http.StatusInternalServerError, Status: "INTERNAL"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsRateLimitError(tt.err))
		})
	}
}
