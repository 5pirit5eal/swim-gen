package pdf

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
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
