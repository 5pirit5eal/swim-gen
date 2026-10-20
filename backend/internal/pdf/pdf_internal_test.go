package pdf

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/stretchr/testify/assert"
)

func TestContentWithEquipment(t *testing.T) {
	tests := []struct {
		name string
		row  models.Row
		lang models.Language
		want string
	}{
		{
			name: "keeps content without equipment",
			row:  models.Row{Content: "Easy freestyle"},
			lang: models.LanguageEN,
			want: "Easy freestyle",
		},
		{
			name: "appends localized English equipment",
			row: models.Row{
				Content:   "Kick set",
				Equipment: []models.EquipmentType{models.EquipmentFins, models.EquipmentKickboard},
			},
			lang: models.LanguageEN,
			want: "Kick set | Fins, Kickboard",
		},
		{
			name: "appends German equipment",
			row: models.Row{
				Content:   "Technik",
				Equipment: []models.EquipmentType{models.EquipmentPaddles, models.EquipmentSnorkel},
			},
			lang: models.LanguageDE,
			want: "Technik | Handpaddles, Schnorchel",
		},
		{
			name: "renders equipment without empty separator",
			row: models.Row{
				Equipment: []models.EquipmentType{models.EquipmentBuoy},
			},
			lang: models.LanguageEN,
			want: "Pull buoy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contentWithEquipment(tt.row, tt.lang))
		})
	}
}

func TestLargeFontColumnWidths(t *testing.T) {
	standard := getColumnWidths(false)
	large := getColumnWidths(true)

	assert.Equal(t, standard.amount*3, large.amount)
	assert.Equal(t, standard.multiplier*3, large.multiplier)
	assert.Equal(t, standard.distance*3, large.distance)
	assert.Equal(t, 43, large.description)
	assert.Equal(t, 100, large.amount+large.multiplier+large.distance+large.breakTime+large.description+large.intensity+large.volume)
}

func extractRowColumnValues(r interface{ GetColumns() []core.Col }) []string {
	cols := r.GetColumns()
	values := make([]string, len(cols))
	for i, c := range cols {
		children := c.GetStructure().GetNexts()
		if len(children) > 0 {
			if val, ok := children[0].GetData().Value.(string); ok {
				values[i] = val
			}
		}
	}
	return values
}

func TestSubRowAmountAndIntensityRendering(t *testing.T) {
	subRows := []models.Row{
		{
			Amount:     1,
			Multiplier: "x",
			Distance:   800,
			Break:      "30",
			Content:    "Kraul locker",
			Intensity:  "GA2", // Same as parent -> should be empty
			Sum:        800,
		},
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   150,
			Break:      "20",
			Content:    "Kraul mit Flossen",
			Intensity:  "GA1", // Differs from parent -> should be "GA1"
			Sum:        600,
		},
		{
			Amount:     10,
			Multiplier: "", // Empty multiplier -> should default to "x" when Amount > 1
			Distance:   50,
			Break:      "0", // Pause is "0" -> should be omitted
			Content:    "Sprint Kraul",
			Intensity:  "1:35/100m", // Tempo value differing from parent -> should be shown
			Sum:        500,
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   100,
			Break:      "   ", // Whitespace break -> should be empty
			Content:    "Technik",
			Intensity:  "", // Empty intensity -> should be empty
			Sum:        200,
		},
	}

	rendered := createSubRows(subRows, "GA2", props.Text{}, models.LanguageDE, "", getColumnWidths(false), 1)
	assert.Len(t, rendered, 4)

	// Subrow 1: Amount == 1 (omitted), Break "30" (shown), Intensity == parent "GA2" (omitted)
	assert.Equal(t, []string{"", "", "800", "30", "Kraul locker", "", ""}, extractRowColumnValues(rendered[0]))

	// Subrow 2: Amount == 4 ("4", multiplier "x"), Break "20" (shown), Intensity "GA1" != parent "GA2" (shown)
	assert.Equal(t, []string{"4", "x", "150", "20", "Kraul mit Flossen", "GA1", ""}, extractRowColumnValues(rendered[1]))

	// Subrow 3: Amount == 10 ("10", default multiplier "x"), Break "0" (omitted), Tempo/Intensity "1:35/100m" != parent "GA2" (shown)
	assert.Equal(t, []string{"10", "x", "50", "", "Sprint Kraul", "1:35/100m", ""}, extractRowColumnValues(rendered[2]))

	// Subrow 4: Amount == 2 ("2", multiplier "x"), whitespace Break (omitted), empty Intensity (omitted)
	assert.Equal(t, []string{"2", "x", "100", "", "Technik", "", ""}, extractRowColumnValues(rendered[3]))

	// Verify visual cues: left accent border on subrow cells and italic muted text on subrow components
	for _, r := range rendered {
		rowDetails := r.GetStructure().GetData().Details
		assert.Equal(t, border.Left, rowDetails["prop_border_type"])
		assert.Equal(t, (&props.Color{Red: 80, Green: 80, Blue: 80}).ToString(), rowDetails["prop_border_color"])

		textDetails := r.GetColumns()[2].GetStructure().GetNexts()[0].GetData().Details
		assert.Equal(t, fontstyle.Italic, textDetails["prop_font_style"])
	}
}
