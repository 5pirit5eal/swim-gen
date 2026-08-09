package models_test

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/require"
)

func TestCalculateCSSPace(t *testing.T) {
	time200, time400 := 200, 420

	pace, ok := models.CalculateCSSPace(&time200, &time400)

	require.True(t, ok)
	require.Equal(t, 110.0, pace)
}

func TestCalculateCSSPaceRejectsInvalidTimes(t *testing.T) {
	time200, time400 := 420, 200

	_, ok := models.CalculateCSSPace(&time200, &time400)

	require.False(t, ok)
}

func TestCalculateCSSZones(t *testing.T) {
	zones := models.CalculateCSSZones(110)

	require.Len(t, zones, 5)
	require.Equal(t, "Schwelle/LT", zones[2].Name)
	require.Equal(t, "Anaerob/VO2max", zones[3].Name)
	require.Equal(t, "Sprint", zones[4].Name)
	require.Equal(t, "2:06", models.FormatPace(zones[0].FasterPaceSeconds))
	require.Equal(t, "2:23", models.FormatPace(zones[0].SlowerPaceSeconds))
	require.Nil(t, zones[4].MaxSpeedPercent)
}
