package garmin_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupPlan() models.Plan {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   100,
			Break:      "30",
			Content:    "Kraul-Beine",
			Intensity:  "GA1",
			Sum:        200,
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling",
			Intensity:  "TÜ",
			Sum:        100,
		},
		{
			Amount:     0,
			Multiplier: "",
			Distance:   0,
			Break:      "",
			Content:    "Gesamt",
			Intensity:  "",
			Sum:        300,
		},
	}
	return models.Plan{
		Table:       table,
		Title:       "test plan",
		Description: "test description",
		PlanID:      "test123",
	}
}

func TestConvertTrainingPlanToFit(t *testing.T) {
	plan := setupPlan()

	// Call UpdateSum to recalculate the sums
	fit, err := garmin.ConvertTrainingPlanToFit(&plan, 25)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, fit)
}

func TestEncodeFITFile(t *testing.T) {
	plan := setupPlan()
	fit, err := garmin.ConvertTrainingPlanToFit(&plan, 25)
	assert.NoError(t, err)
	assert.NotNil(t, fit)

	var buf bytes.Buffer
	err = garmin.EncodeFITFile(context.Background(), &buf, fit)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf.Bytes())
}
