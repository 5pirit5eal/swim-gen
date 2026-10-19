package pdf

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
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
			Amount:     2,
			Multiplier: "x",
			Distance:   150,
			Break:      "20",
			Content:    "Kraul mit Flossen",
			Intensity:  "GA1", // Differs from parent -> should be "GA1"
			Sum:        300,
		},
		{
			Amount:     3,
			Multiplier: "", // Empty multiplier -> should default to "x" when Amount > 1
			Distance:   50,
			Break:      "",
			Content:    "Sprint Kraul",
			Intensity:  "", // Empty intensity -> should be empty
			Sum:        150,
		},
	}

	rendered := createSubRows(subRows, "GA2", props.Text{}, models.LanguageDE, "", getColumnWidths(false), 1)
	assert.Len(t, rendered, 3)

	// Subrow 1: Amount == 1 (omitted), Intensity == parent "GA2" (omitted)
	assert.Equal(t, []string{"", "", "800", "", "Kraul locker", "", ""}, extractRowColumnValues(rendered[0]))

	// Subrow 2: Amount == 2 > 1 (shown with multiplier), Intensity "GA1" != parent "GA2" (shown)
	assert.Equal(t, []string{"2", "x", "150", "", "Kraul mit Flossen", "GA1", ""}, extractRowColumnValues(rendered[1]))

	// Subrow 3: Amount == 3 > 1 with empty multiplier (defaults to "x"), empty intensity (omitted)
	assert.Equal(t, []string{"3", "x", "50", "", "Sprint Kraul", "", ""}, extractRowColumnValues(rendered[2]))
}
