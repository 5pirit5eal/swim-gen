package genai

import "testing"

func TestFormatQueryEmbeddingInput(t *testing.T) {
	content := "four hundred meter endurance\n\nBenutzerprofil Präferenzen:"
	want := "task: search result | query: four hundred meter endurance\n\nBenutzerprofil Präferenzen:"

	if got := formatQueryEmbeddingInput(content); got != want {
		t.Fatalf("formatQueryEmbeddingInput() = %q, want %q", got, want)
	}
}
