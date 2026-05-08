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
	pdfBytes, err := pdf.GenerateEasyReadablePDF(&table, false, models.LanguageDE, "")

	// Assertions
	assert.NoError(t, err, "TableToPDF should not return an error")
	assert.NotEmpty(t, pdfBytes, "TableToPDF should return non-empty PDF bytes")

	// Test that the bytes are writable to pdf file
	err = writePDF("table.pdf", pdfBytes)

	assert.NoError(t, err, "writePDF should not return an error")

	// Cleanup
	err = os.Remove("table.pdf")
	assert.NoError(t, err, "Cleanup failed")

	// Test large font
	pdfBytes, err = pdf.GenerateEasyReadablePDF(&table, true, models.LanguageDE, "")
	assert.NoError(t, err, "TableToPDF with large font should not return an error")
	assert.NotEmpty(t, pdfBytes, "TableToPDF with large font should return non-empty PDF bytes")

	err = writePDF("table_lf.pdf", pdfBytes)
	assert.NoError(t, err, "writePDF with large font should not return an error")

	err = os.Remove("table_lf.pdf")
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

	planPDF, err := pdf.PlanToPDF(plan, false, false, models.LanguageDE, "")
	assert.NoError(t, err, "PlanToPDF should not return an error")
	assert.NotEmpty(t, planPDF, "PlanToPDF should return non-empty PDF bytes")

	// Test that the bytes are writable to pdf file
	err = writePDF("plan.pdf", planPDF)

	assert.NoError(t, err, "writePDF should not return an error")

	// Cleanup
	err = os.Remove("plan.pdf")
	assert.NoError(t, err, "Cleanup failed")
}

func TestGenerateStoragePath(t *testing.T) {
	tests := []struct {
		name     string
		username string
		planID   string
		title    string
		want     string
	}{
		{
			name:     "Authenticated user with title",
			username: "johndoe",
			planID:   "plan123",
			title:    "My Training Plan",
			want:     "johndoe/my_training_plan.pdf",
		},
		{
			name:     "Authenticated user with special chars in title",
			username: "johndoe",
			planID:   "plan123",
			title:    "Technik-Tüftler & Ausdauer-As",
			want:     "johndoe/technik-tueftler_ausdauer-as.pdf",
		},
		{
			name:     "Authenticated user empty title",
			username: "johndoe",
			planID:   "plan123",
			title:    "",
			want:     "johndoe/training-plan.pdf",
		},
		{
			name:     "Anonymous with PlanID",
			username: "",
			planID:   "plan123",
			title:    "My Plan",
			want:     "plan123/my_plan.pdf",
		},
		{
			name:     "Anonymous with PlanID empty title",
			username: "",
			planID:   "plan123",
			title:    "",
			want:     "plan123/training-plan.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pdf.GenerateStoragePath(tt.username, tt.planID, tt.title)
			assert.Equal(t, tt.want, got)
		})
	}

	t.Run("Fallback to anonymous uuid", func(t *testing.T) {
		got := pdf.GenerateStoragePath("", "", "Some Title")
		assert.Contains(t, got, "anonymous/")
		assert.True(t, len(got) > len("anonymous/.pdf"))
	})
}

// TestHyperlinks creates PDFs with various hyperlink scenarios for visual inspection.
// The PDFs are saved to the current directory and are cleaned up by default; they are only
// retained for review when the GENERATE_PDF environment variable is set.
func TestHyperlinks(t *testing.T) {
	baseURL := "https://example.com"

	// Test table with various hyperlink scenarios
	table := models.Table{
		// Row 1: Single hyperlink
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "See [this drill](https://example.com/drill1) for technique",
			Intensity:  "GA1",
			Sum:        100,
		},
		// Row 2: Multiple hyperlinks in one cell
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   50,
			Break:      "30",
			Content:    "Practice [drill A](/drills/a) and [drill B](/drills/b) alternating",
			Intensity:  "TÜ",
			Sum:        100,
		},
		// Row 3: Multi-line text BEFORE a hyperlink (tests overlap issue)
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   75,
			Break:      "20",
			Content:    "This is a very long instruction that should wrap to multiple lines in the PDF and then includes a hyperlink at the end [click here](/info)",
			Intensity:  "GA1",
			Sum:        300,
		},
		// Row 4: Hyperlink in brackets (tests formatting issue)
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   200,
			Break:      "",
			Content:    "Warm up slowly (see [video tutorial](/video) for guidance)",
			Intensity:  "ReKom",
			Sum:        200,
		},
		// Row 5: URL at end of content block (tests end-of-block bug)
		{
			Amount:     3,
			Multiplier: "x",
			Distance:   100,
			Break:      "45",
			Content:    "Sprint section [details](/sprint)",
			Intensity:  "S",
			Sum:        300,
		},
		// Row 6: Text after hyperlink (tests capturing text after last match)
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   150,
			Break:      "60",
			Content:    "Check [form guide](/form) and maintain proper technique throughout",
			Intensity:  "GA2",
			Sum:        300,
		},
		// Row 7: No hyperlinks (baseline test)
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   100,
			Break:      "",
			Content:    "Cool down with easy swimming",
			Intensity:  "ReKom",
			Sum:        100,
		},
		// Row 8: Multi-line with multiple hyperlinks (stress test)
		{
			Amount:     5,
			Multiplier: "x",
			Distance:   50,
			Break:      "15",
			Content:    "This exercise combines [technique A](/tech-a) with [technique B](/tech-b) for optimal results. Make sure to review both beforehand.",
			Intensity:  "TÜ",
			Sum:        250,
		},
	}
	table.AddSum()
	table.UpdateSum()

	tests := []struct {
		name       string
		horizontal bool
		largeFont  bool
		filename   string
	}{
		{
			name:       "Standard PDF",
			horizontal: false,
			largeFont:  false,
			filename:   "test_hyperlinks_standard.pdf",
		},
		{
			name:       "Large font (easy-to-read) PDF",
			horizontal: false,
			largeFont:  true,
			filename:   "test_hyperlinks_largefont.pdf",
		},
		{
			name:       "Horizontal PDF",
			horizontal: true,
			largeFont:  false,
			filename:   "test_hyperlinks_horizontal.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pdfBytes []byte
			var err error

			if tt.largeFont {
				pdfBytes, err = pdf.GenerateEasyReadablePDF(&table, tt.horizontal, models.LanguageEN, baseURL)
			} else {
				plan := &models.Plan{
					Title:       "Hyperlink Test Plan",
					Description: "Testing various hyperlink scenarios in PDF generation",
					Table:       table,
				}
				pdfBytes, err = pdf.PlanToPDF(plan, tt.horizontal, false, models.LanguageEN, baseURL)
			}

			if err != nil {
				t.Fatalf("Failed to generate PDF: %v", err)
			}

			if len(pdfBytes) == 0 {
				t.Fatal("Generated PDF is empty")
			}

			// Write PDF for visual inspection (not cleaned up)
			err = writePDF(tt.filename, pdfBytes)
			if err != nil {
				t.Fatalf("Failed to write PDF: %v", err)
			}

			// Cleanup
			if os.Getenv("GENERATE_PDF") == "" {
				err = os.Remove(tt.filename)
				if err != nil {
					t.Fatalf("Failed to cleanup PDF: %v", err)
				}
			}

			t.Logf("Generated %s - please inspect visually", tt.filename)
		})
	}
}

// TestSubRows creates PDFs with subrows (compound sets) for visual inspection.
// Subrows should have:
// - Empty amount field
// - Distance value shown in distance column
// - Empty sum field
// - Visual indent with "↳" prefix
// - Alternating background colors
func TestSubRows(t *testing.T) {
	baseURL := "https://example.com"

	// Test table with subrows (compound sets like 8 x (800 + 200))
	table := models.Table{
		// Row with subrows - main set
		{
			Amount:     8,
			Multiplier: "x",
			Distance:   1000, // Will be calculated from subrows
			Break:      "60",
			Content:    "Hauptset - Ausdauer",
			Intensity:  "GA2",
			Sum:        8000,
			SubRows: []models.Row{
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   800,
					Break:      "30",
					Content:    "Kraul locker",
					Intensity:  "GA1",
					Sum:        800,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "",
					Content:    "Rücken Technik",
					Intensity:  "TÜ",
					Sum:        200,
				},
			},
		},
		// Another row with subrows
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   400,
			Break:      "45",
			Content:    "Sprint-Vorbereitung",
			Intensity:  "GA2",
			Sum:        1600,
			SubRows: []models.Row{
				{
					Amount:     2,
					Multiplier: "x",
					Distance:   150,
					Break:      "20",
					Content:    "Kraul mit Flossen",
					Intensity:  "GA2",
					Sum:        300,
				},
				{
					Amount:     2,
					Multiplier: "x",
					Distance:   50,
					Break:      "",
					Content:    "Sprint Kraul",
					Intensity:  "S",
					Sum:        100,
				},
			},
		},
		// Regular row without subrows
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   400,
			Break:      "",
			Content:    "Locker ausschwimmen",
			Intensity:  "ReKom",
			Sum:        400,
		},
		// Row with nested subrows (depth 2)
		{
			Amount:     3,
			Multiplier: "x",
			Distance:   600,
			Break:      "90",
			Content:    "Technik-Block",
			Intensity:  "TÜ",
			Sum:        1800,
			SubRows: []models.Row{
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "30",
					Content:    "Sculling",
					Intensity:  "TÜ",
					Sum:        200,
					SubRows: []models.Row{
						{
							Amount:   1,
							Distance: 100,
							Content:  "Einarmig",
							Sum:      100,
						},
						{
							Amount:   1,
							Distance: 100,
							Content:  "Zweifarmer",
							Sum:      100,
						},
					},
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "30",
					Content:    "Seitlage",
					Intensity:  "TÜ",
					Sum:        200,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "",
					Content:    "Atemtechnik",
					Intensity:  "TÜ",
					Sum:        200,
				},
			},
		},
	}
	table.AddSum()
	table.UpdateSum()

	tests := []struct {
		name       string
		horizontal bool
		largeFont  bool
		filename   string
	}{
		{
			name:       "SubRows Standard PDF",
			horizontal: false,
			largeFont:  false,
			filename:   "test_subrows_standard.pdf",
		},
		{
			name:       "SubRows Large Font PDF",
			horizontal: false,
			largeFont:  true,
			filename:   "test_subrows_largefont.pdf",
		},
		{
			name:       "SubRows Horizontal PDF",
			horizontal: true,
			largeFont:  false,
			filename:   "test_subrows_horizontal.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pdfBytes []byte
			var err error

			if tt.largeFont {
				pdfBytes, err = pdf.GenerateEasyReadablePDF(&table, tt.horizontal, models.LanguageDE, baseURL)
			} else {
				plan := &models.Plan{
					Title:       "SubRows Test Plan",
					Description: "Testing subrow rendering in PDF generation",
					Table:       table,
				}
				pdfBytes, err = pdf.PlanToPDF(plan, tt.horizontal, false, models.LanguageDE, baseURL)
			}

			if err != nil {
				t.Fatalf("Failed to generate PDF: %v", err)
			}

			if len(pdfBytes) == 0 {
				t.Fatal("Generated PDF is empty")
			}

			// Write PDF for visual inspection
			err = writePDF(tt.filename, pdfBytes)
			if err != nil {
				t.Fatalf("Failed to write PDF: %v", err)
			}

			// Cleanup unless GENERATE_PDF is set
			if os.Getenv("GENERATE_PDF") == "" {
				err = os.Remove(tt.filename)
				if err != nil {
					t.Fatalf("Failed to cleanup PDF: %v", err)
				}
			}

			t.Logf("Generated %s - please inspect visually for proper subrow indentation", tt.filename)
		})
	}
}

// TestSubRowsWithHyperlinks tests subrows that contain markdown links
func TestSubRowsWithHyperlinks(t *testing.T) {
	baseURL := "https://example.com"

	table := models.Table{
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   300,
			Break:      "45",
			Content:    "Technik-Set mit Links [Anleitung](/technik)",
			Intensity:  "TÜ",
			Sum:        1200,
			SubRows: []models.Row{
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   100,
					Break:      "20",
					Content:    "Kraul [Video](/video1)",
					Intensity:  "GA1",
					Sum:        100,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   100,
					Break:      "20",
					Content:    "Rücken [Tutorial](/tutorial)",
					Intensity:  "GA1",
					Sum:        100,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   100,
					Break:      "",
					Content:    "Brust [Guide](/guide)",
					Intensity:  "GA1",
					Sum:        100,
				},
			},
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

	pdfBytes, err := pdf.GenerateEasyReadablePDF(&table, true, models.LanguageDE, baseURL)
	assert.NoError(t, err, "PDF generation with subrows and hyperlinks should not fail")
	assert.NotEmpty(t, pdfBytes, "Generated PDF should not be empty")

	err = writePDF("test_subrows_hyperlinks.pdf", pdfBytes)
	assert.NoError(t, err, "Writing PDF should not fail")

	if os.Getenv("GENERATE_PDF") == "" {
		err = os.Remove("test_subrows_hyperlinks.pdf")
		assert.NoError(t, err, "Cleanup should succeed")
	}
}

// TestDeeplyNestedSubRows tests 3+ levels of nested subrows with content aggregation
func TestDeeplyNestedSubRows(t *testing.T) {
	baseURL := "https://example.com"

	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   800,
			Break:      "60",
			Content:    "Hauptset",
			Intensity:  "GA2",
			Sum:        1600,
			SubRows: []models.Row{
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   400,
					Break:      "30",
					Content:    "Ausdauer",
					Intensity:  "GA1",
					Sum:        400,
					SubRows: []models.Row{
						{
							Amount:   1,
							Distance: 200,
							Content:  "Kraul",
							Sum:      200,
							SubRows: []models.Row{
								{
									Amount:   1,
									Distance: 100,
									Content:  "locker",
									Sum:      100,
								},
								{
									Amount:   1,
									Distance: 100,
									Content:  "technisch",
									Sum:      100,
								},
							},
						},
						{
							Amount:   1,
							Distance: 200,
							Content:  "Rücken",
							Sum:      200,
						},
					},
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   400,
					Break:      "",
					Content:    "Sprint",
					Intensity:  "S",
					Sum:        400,
				},
			},
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

	// Test with large font (easy-to-read)
	pdfBytes, err := pdf.GenerateEasyReadablePDF(&table, true, models.LanguageDE, baseURL)
	assert.NoError(t, err, "PDF generation with deeply nested subrows should not fail")
	assert.NotEmpty(t, pdfBytes, "Generated PDF should not be empty")

	err = writePDF("test_deeply_nested.pdf", pdfBytes)
	assert.NoError(t, err, "Writing PDF should not fail")

	if os.Getenv("GENERATE_PDF") == "" {
		err = os.Remove("test_deeply_nested.pdf")
		assert.NoError(t, err, "Cleanup should succeed")
	}

	t.Log("Generated test_deeply_nested.pdf - please verify aggregated content format:")
	t.Log("Expected: ↳ Ausdauer (200m Kraul (100m locker + 100m technisch) + 200m Rücken)")
}
