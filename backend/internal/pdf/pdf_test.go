package pdf_test

import (
	"os"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/5pirit5eal/swim-gen/internal/pdf"
	"github.com/stretchr/testify/assert"
)

func writePDF(filename string, pdfBytes []byte) error {
	return os.WriteFile(filename, pdfBytes, 0644)

}

func TestTableToPDF(t *testing.T) {
	// Create a sample table
	table := models.Table{
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   200,
			Break:      "",
			Content:    "Einschwimmen",
			Intensity:  "",
			Sum:        200,
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   100,
			Break:      "30",
			Content:    "Kraul-Beine m. Kurzflossen + Schnorchel jeweils 50m Streamline + 50m Schultern an der Wasseroberfläche halten",
			Intensity:  "GA1",
			Sum:        200,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling mit Schnorchel Beinarbeit-Timing beachten",
			Intensity:  "TÜ",
			Sum:        200,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   50,
			Break:      "30",
			Content:    "Kraul Flossen, Paddles, Schnorchel 3 Del-Kicks Unterwasser + Züge zählen und „distance per stroke“",
			Intensity:  "TÜ",
			Sum:        200,
		},
		{
			Amount:     6,
			Multiplier: "x",
			Distance:   30,
			Break:      "30",
			Content:    "15m Kraul-WASSER-Start „Bursts“ + 15m lo.",
			Intensity:  "S",
			Sum:        180,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "Locker schwimmen als aktive Pause",
			Intensity:  "ReKom",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "60",
			Content:    "15m Spurt Breakout + 85m locker",
			Intensity:  "S",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "120",
			Content:    "25m „easy-Speed-95%“ + 75m locker",
			Intensity:  "S",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "180",
			Content:    "35m Spurt Tempoaufbau + 65m locker",
			Intensity:  "S",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "50m Spurt „alle Punkte umsetzen“ + 50m locker schwimmen als aktive Pause",
			Intensity:  "S/WA",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   400,
			Break:      "",
			Content:    "Locker beliebig mit Kurzflossen",
			Intensity:  "ReKom",
			Sum:        400,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "Kraul/Rücken-Beine",
			Intensity:  "ReKom/GA1",
			Sum:        400,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   200,
			Break:      "",
			Content:    "Ausschwimmen",
			Intensity:  "ReKom",
			Sum:        200,
		},
	}
	table.AddSum()
	table.UpdateSum()

	// Call the function to be tested
	pdfBytes, err := pdf.TableToPDF(table)

	// Assertions
	assert.NoError(t, err, "TableToPDF should not return an error")
	assert.NotEmpty(t, pdfBytes, "TableToPDF should return non-empty PDF bytes")

	// Test that the bytes are writable to pdf file
	err = writePDF("table.pdf", pdfBytes)

	assert.NoError(t, err, "writePDF should not return an error")

	// Cleanup
	err = os.Remove("table.pdf")
	assert.NoError(t, err, "Cleanup failed")
}

func TestPlantoPDF(t *testing.T) {
	table := models.Table{
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   200,
			Break:      "",
			Content:    "Einschwimmen",
			Intensity:  "",
			Sum:        200,
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   100,
			Break:      "30",
			Content:    "Kraul-Beine m. Kurzflossen + Schnorchel jeweils 50m Streamline + 50m Schultern an der Wasseroberfläche halten",
			Intensity:  "GA1",
			Sum:        200,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling mit Schnorchel Beinarbeit-Timing beachten",
			Intensity:  "TÜ",
			Sum:        200,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   50,
			Break:      "30",
			Content:    "Kraul Flossen, Paddles, Schnorchel 3 Del-Kicks Unterwasser + Züge zählen und „distance per stroke“",
			Intensity:  "TÜ",
			Sum:        200,
		},
		{
			Amount:     6,
			Multiplier: "x",
			Distance:   30,
			Break:      "30",
			Content:    "15m Kraul-WASSER-Start „Bursts“ + 15m lo.",
			Intensity:  "S",
			Sum:        180,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "Locker schwimmen als aktive Pause",
			Intensity:  "ReKom",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "60",
			Content:    "15m Spurt Breakout + 85m locker",
			Intensity:  "S",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "120",
			Content:    "25m „easy-Speed-95%“ + 75m locker",
			Intensity:  "S",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "180",
			Content:    "35m Spurt Tempoaufbau + 65m locker",
			Intensity:  "S",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "50m Spurt „alle Punkte umsetzen“ + 50m locker schwimmen als aktive Pause",
			Intensity:  "S/WA",
			Sum:        100,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   400,
			Break:      "",
			Content:    "Locker beliebig mit Kurzflossen",
			Intensity:  "ReKom",
			Sum:        400,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "Kraul/Rücken-Beine",
			Intensity:  "ReKom/GA1",
			Sum:        400,
		},
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   200,
			Break:      "",
			Content:    "Ausschwimmen",
			Intensity:  "ReKom",
			Sum:        200,
		},
	}
	// Create a sample plan
	plan := &models.Plan{
		Title: "Kraul-Sprint Training für Höchstgeschwindigkeit",
		Description: `Dieser Trainingsplan ist ein Super-Sprint-Plan, inspiriert von einem Olympiasieger. 
Er konzentriert sich auf die Entwicklung deiner absoluten Höchstgeschwindigkeit im Kraulschwimmen durch kurze, 
intensive Sprint-Abschnitte und spezifische Technikübungen. Achte auf die Einhaltung der Pausen, 
um dich optimal auf die Sprints vorzubereiten.`,
		Table: table,
	}

	planPDF, err := pdf.PlanToPDF(plan)
	assert.NoError(t, err, "PlanToPDF should not return an error")
	assert.NotEmpty(t, planPDF, "PlanToPDF should return non-empty PDF bytes")

	// Test that the bytes are writable to pdf file
	err = writePDF("plan.pdf", planPDF)

	assert.NoError(t, err, "writePDF should not return an error")

	// Cleanup
	err = os.Remove("plan.pdf")
	assert.NoError(t, err, "Cleanup failed")
}
