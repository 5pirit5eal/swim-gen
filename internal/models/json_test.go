package models_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/5pirit5eal/swim-rag/internal/models"
)

type Example struct {
	Title       string `json:"title"`
	Description int    `json:"description"`
}

type ExampleList []Example

func TestStaticToStruct(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		Expected any
	}{
		{
			name: "Valid JSON",
			data: `{"title": "Test Title", "description": 1}`,
			Expected: Example{
				Title:       "Test Title",
				Description: 1,
			},
		},
		{
			name:     "Invalid JSON",
			data:     `{"title": "Test Title", "description": "invalid"}`,
			Expected: nil,
		},
		{
			name: "List of JSON",
			data: `[{"title": "Test Title", "description": 1}, {"title": "Another Title", "description": 2}]`,
			Expected: []Example{
				{
					Title:       "Test Title",
					Description: 1,
				},
				{
					Title:       "Another Title",
					Description: 2,
				},
			},
		},
		{
			name: "List of JSON with custom type",
			data: `[{"title": "Test Title", "description": 1}, {"title": "Another Title", "description": 2}]`,
			Expected: ExampleList{
				{
					Title:       "Test Title",
					Description: 1,
				},
				{
					Title:       "Another Title",
					Description: 2,
				},
			},
		},
		{
			name: "Table JSON",
			data: `[{"Amount":1,"Multiplier":"x","Distance":300,"Break":"","Content":"Einschwimmen","Intensity":"","Sum":300},{"Amount":1,"Multiplier":"x","Distance":200,"Break":"30","Content":"Kraul-Armarbeit mit Pullbuoy","Intensity":"GA1","Sum":200},{"Amount":1,"Multiplier":"x","Distance":200,"Break":"30","Content":"Beinarbeit mit Brett","Intensity":"GA1","Sum":200},{"Amount":1,"Multiplier":"x","Distance":200,"Break":"30","Content":"Lagenwechsel, beliebige Aufteilung","Intensity":"GA1","Sum":200},{"Amount":4,"Multiplier":"x","Distance":100,"Break":"60","Content":"15m Spurt HSA (Startblock) + 85m ReKom","Intensity":"S","Sum":400},{"Amount":1,"Multiplier":"x","Distance":100,"Break":"","Content":"Locker schwimmen als aktive Pause","Intensity":"","Sum":100},{"Amount":4,"Multiplier":"x","Distance":100,"Break":"60","Content":"25m Spurt HSA (Startblock) + 75m ReKom","Intensity":"S","Sum":400},{"Amount":1,"Multiplier":"x","Distance":100,"Break":"","Content":"Locker schwimmen als aktive Pause","Intensity":"","Sum":100},{"Amount":4,"Multiplier":"x","Distance":200,"Break":"60","Content":"50m Spurt HSA (Startblock) + 150m ReKom","Intensity":"S/SA","Sum":800},{"Amount":1,"Multiplier":"x","Distance":300,"Break":"","Content":"Ausschwimmen","Intensity":"ReKom","Sum":300},{"Amount":0,"Multiplier":"","Distance":0,"Break":"","Content":"Gesamt","Intensity":"","Sum":3000}]`,
			Expected: models.Table{
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   300,
					Break:      "",
					Content:    "Einschwimmen",
					Intensity:  "",
					Sum:        300,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "30",
					Content:    "Kraul-Armarbeit mit Pullbuoy",
					Intensity:  "GA1",
					Sum:        200,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "30",
					Content:    "Beinarbeit mit Brett",
					Intensity:  "GA1",
					Sum:        200,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   200,
					Break:      "30",
					Content:    "Lagenwechsel, beliebige Aufteilung",
					Intensity:  "GA1",
					Sum:        200,
				},
				{
					Amount:     4,
					Multiplier: "x",
					Distance:   100,
					Break:      "60",
					Content:    "15m Spurt HSA (Startblock) + 85m ReKom",
					Intensity:  "S",
					Sum:        400,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   100,
					Break:      "",
					Content:    "Locker schwimmen als aktive Pause",
					Intensity:  "",
					Sum:        100,
				},
				{
					Amount:     4,
					Multiplier: "x",
					Distance:   100,
					Break:      "60",
					Content:    "25m Spurt HSA (Startblock) + 75m ReKom",
					Intensity:  "S",
					Sum:        400,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   100,
					Break:      "",
					Content:    "Locker schwimmen als aktive Pause",
					Intensity:  "",
					Sum:        100,
				},
				{
					Amount:     4,
					Multiplier: "x",
					Distance:   200,
					Break:      "60",
					Content:    "50m Spurt HSA (Startblock) + 150m ReKom",
					Intensity:  "S/SA",
					Sum:        800,
				},
				{
					Amount:     1,
					Multiplier: "x",
					Distance:   300,
					Break:      "",
					Content:    "Ausschwimmen",
					Intensity:  "ReKom",
					Sum:        300,
				},
				{
					Amount:     0,
					Multiplier: "",
					Distance:   0,
					Break:      "",
					Content:    "Gesamt",
					Intensity:  "",
					Sum:        3000,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var data interface{}

			err := json.Unmarshal([]byte(tc.data), &data)
			if err != nil {
				t.Errorf("Failed to unmarshal JSON: %v", err)
			}
			var example = reflect.New(reflect.TypeOf(tc.Expected)).Elem().Interface()
			err = models.JSONInterfaceToStruct(data, example)
			if err != nil && tc.Expected != nil {
				t.Errorf("Failed to convert static to struct: %v", err)
			}
			if tc.Expected == nil {
				return
			}
			if !reflect.DeepEqual(tc.Expected, reflect.Indirect(reflect.ValueOf(example)).Interface()) {
				t.Errorf("Expected %v, got %v", tc.Expected, example)
			}
		})
	}
}
