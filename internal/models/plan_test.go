package models_test

import (
	"testing"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSum(t *testing.T) {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   100,
			Break:      "30",
			Content:    "Kraul-Beine",
			Intensity:  "GA1",
			Sum:        0, // Initial sum is incorrect
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling",
			Intensity:  "TÜ",
			Sum:        0, // Initial sum is incorrect
		},
		{
			Amount:     0,
			Multiplier: "",
			Distance:   0,
			Break:      "",
			Content:    "Gesamt",
			Intensity:  "",
			Sum:        0, // Initial sum is incorrect
		},
	}

	// Call UpdateSum to recalculate the sums
	table.UpdateSum()

	// Assertions
	assert.Equal(t, 200, table[0].Sum, "Sum for the first row should be 200")
	assert.Equal(t, 100, table[1].Sum, "Sum for the second row should be 100")
	assert.Equal(t, 300, table[2].Sum, "Sum for the third row should be 300")
}
