package models

import "fmt"

type CSSZone struct {
	Name              string
	Focus             string
	MinSpeedPercent   int
	MaxSpeedPercent   *int
	FasterPaceSeconds float64
	SlowerPaceSeconds float64
}

func CalculateCSSPace(time200m, time400m *int) (float64, bool) {
	if time200m == nil || time400m == nil || *time200m <= 0 || *time400m <= *time200m {
		return 0, false
	}

	return float64(*time400m-*time200m) / 2, true
}

func CalculateCSSZones(cssPace float64) []CSSZone {
	if cssPace <= 0 {
		return nil
	}

	definitions := []struct {
		name  string
		focus string
		min   int
		max   *int
	}{
		{"Regeneration/ReKom", "Aufwärmen, aktive Erholung und Technik", 77, intPointer(87)},
		{"Aerobe Ausdauer/GA1", "Lange, gleichmäßige Serien und Grundlagenausdauer", 87, intPointer(94)},
		{"Schwelle/LT", "Verbesserung der nachhaltig haltbaren Pace", 95, intPointer(104)},
		{"Anaerob/VO2max", "Geschwindigkeit und Laktattoleranz", 104, intPointer(111)},
		{"Sprint", "Maximale Geschwindigkeit und Kraft", 111, nil},
	}

	zones := make([]CSSZone, 0, len(definitions))
	for _, definition := range definitions {
		maxSpeed := definition.min
		if definition.max != nil {
			maxSpeed = *definition.max
		}
		zones = append(zones, CSSZone{
			Name:              definition.name,
			Focus:             definition.focus,
			MinSpeedPercent:   definition.min,
			MaxSpeedPercent:   definition.max,
			FasterPaceSeconds: cssPace / (float64(maxSpeed) / 100),
			SlowerPaceSeconds: cssPace / (float64(definition.min) / 100),
		})
	}

	return zones
}

func FormatPace(seconds float64) string {
	totalSeconds := int(seconds + 0.5)
	return fmt.Sprintf("%d:%02d", totalSeconds/60, totalSeconds%60)
}

func intPointer(value int) *int {
	return &value
}
