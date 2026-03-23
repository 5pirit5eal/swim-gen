package models_test

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
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
			Sum:        0,
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling",
			Intensity:  "TÜ",
			Sum:        0,
		},
		{
			Amount:     0,
			Multiplier: "",
			Distance:   0,
			Break:      "",
			Content:    "Gesamt",
			Intensity:  "",
			Sum:        0,
		},
	}

	// Call UpdateSum to recalculate the sums
	table.UpdateSum()

	assert.Equal(t, 200, table[0].Sum, "Sum for the first row should be 200")
	assert.Equal(t, 100, table[1].Sum, "Sum for the second row should be 100")
	assert.Equal(t, 300, table[2].Sum, "Sum for the third row should be 300")
}

func TestUpdateSumWithSubRows(t *testing.T) {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   100000, // Wrong distance that should be recalculated based on subRows
			Break:      "30",
			Content:    "Kraul-Beine",
			Intensity:  "GA1",
			Sum:        0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 50, Break: "15", Content: "Freestyle", Intensity: "GA1", Sum: 0},
				{Amount: 1, Distance: 50, Break: "15", Content: "Rücken", Intensity: "GA1", Sum: 0},
			},
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling",
			Intensity:  "TÜ",
			Sum:        0,
		},
		{
			Amount:     0,
			Multiplier: "",
			Distance:   0,
			Break:      "",
			Content:    "Gesamt",
			Intensity:  "",
			Sum:        0,
		},
	}

	// Call UpdateSum to recalculate the sums
	table.UpdateSum()

	assert.Equal(t, 100, table[0].Distance, "Distance for the first row should be 100")
	assert.Equal(t, 200, table[0].Sum, "Sum for the first row should be 200")
	assert.Equal(t, 100, table[1].Sum, "Sum for the second row should be 100")
	assert.Equal(t, 300, table[2].Sum, "Sum for the third row should be 300")
}

func TestGeneratedPlanSchema(t *testing.T) {
	schema, err := models.GeneratedPlanSchema()
	assert.NoError(t, err, "Failed to retrieve schema")

	// check if the schema is valid json
	assert.NotEmpty(t, schema, "Schema should not be empty")
}

func TestUpdateSum_NestedRows(t *testing.T) {
	table := models.Table{
		{
			Amount:     8,
			Multiplier: "x",
			Distance:   0,
			Break:      "20",
			Content:    "Main Set",
			Intensity:  "GA1",
			Sum:        0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 800, Break: "10", Content: "Freestyle", Intensity: "GA1", Sum: 0},
				{Amount: 1, Distance: 200, Break: "0", Content: "Kick", Intensity: "GA1", Sum: 0},
			},
		},
	}
	table.UpdateSum()

	assert.Equal(t, 1000, table[0].Distance, "Parent Distance should be 1000")
	assert.Equal(t, 8000, table[0].Sum, "Parent Sum should be 8000")
	assert.Equal(t, 800, table[0].SubRows[0].Sum, "Child 1 Sum should be recalculated to 800")
	assert.Equal(t, 200, table[0].SubRows[1].Sum, "Child 2 Sum should be recalculated to 200")
}

func TestUpdateSum_NestedRowsWithIncorrectSums(t *testing.T) {
	table := models.Table{
		{
			Amount:     6,
			Multiplier: "x",
			Distance:   0,
			Break:      "15",
			Content:    "Main Set",
			Intensity:  "GA2",
			Sum:        0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 400, Break: "10", Content: "Kraul", Intensity: "GA2", Sum: 9999},
				{Amount: 1, Distance: 100, Break: "5", Content: "Brust", Intensity: "GA2", Sum: 9999},
			},
		},
	}
	table.UpdateSum()

	assert.Equal(t, 500, table[0].Distance, "Main Set Distance should be 500")
	assert.Equal(t, 3000, table[0].Sum, "Main Set sum should be 3000 (recalculated from correct child sums)")
	assert.Equal(t, 400, table[0].SubRows[0].Sum, "Child 1 Sum should be recalculated to 400")
	assert.Equal(t, 100, table[0].SubRows[1].Sum, "Child 2 Sum should be recalculated to 100")
}

func TestUpdateSum_MixedRows(t *testing.T) {
	table := models.Table{
		{Amount: 4, Distance: 100, Break: "20", Content: "Warmup", Intensity: "Rekom", Sum: 400},
		{
			Amount:     6,
			Multiplier: "x",
			Distance:   0,
			Break:      "15",
			Content:    "Main Set",
			Intensity:  "GA2",
			Sum:        0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 400, Break: "10", Content: "Kraul", Intensity: "GA2", Sum: 0},
				{Amount: 1, Distance: 100, Break: "5", Content: "Brust", Intensity: "GA2", Sum: 0},
			},
		},
		{Amount: 1, Distance: 200, Break: "0", Content: "Cooldown", Intensity: "Rekom", Sum: 200},
	}
	table.UpdateSum()

	assert.Equal(t, 400, table[0].Sum, "Warmup sum should be 400")
	assert.Equal(t, 500, table[1].Distance, "Main Set Distance should be 500")
	assert.Equal(t, 3000, table[1].Sum, "Main Set sum should be 3000")
	assert.Equal(t, 200, table[2].Sum, "Cooldown sum should be 200")
}

func TestValidate_ValidTable(t *testing.T) {
	table := models.Table{
		{Amount: 4, Distance: 100, Content: "Test", Intensity: "GA1"},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 100, Content: "Test", Intensity: "GA1"},
			},
		},
	}
	err := table.Validate()
	assert.NoError(t, err, "Valid table should pass validation")
}

func TestValidate_NegativeAmount(t *testing.T) {
	table := models.Table{
		{Amount: -1, Distance: 100, Content: "Test", Intensity: "GA1"},
	}
	err := table.Validate()
	assert.Error(t, err, "Table with negative amount should fail validation")
}

func TestValidate_NegativeDistance(t *testing.T) {
	table := models.Table{
		{Amount: 4, Distance: -100, Content: "Test", Intensity: "GA1"},
	}
	err := table.Validate()
	assert.Error(t, err, "Table with negative distance should fail validation")
}

func TestValidate_ZeroAmountWithSubRows(t *testing.T) {
	table := models.Table{
		{
			Amount:   0,
			Distance: 0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 100, Content: "Test", Intensity: "GA1"},
			},
		},
	}
	err := table.Validate()
	assert.Error(t, err, "Table with zero amount and subRows should fail validation")
}

func TestValidate_MaxDepth(t *testing.T) {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   0,
			SubRows: []models.Row{
				{
					Amount:     2,
					Multiplier: "x",
					Distance:   0,
					SubRows: []models.Row{
						{
							Amount:     2,
							Multiplier: "x",
							Distance:   0,
							SubRows: []models.Row{
								{
									Amount:     2,
									Multiplier: "x",
									Distance:   0,
									SubRows: []models.Row{
										{
											Amount:     2,
											Multiplier: "x",
											Distance:   0,
											SubRows: []models.Row{
												{Amount: 1, Distance: 100, Content: "Test", Intensity: "GA1"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := table.Validate()
	assert.Error(t, err, "Table exceeding max depth (5) should fail validation")
}

func TestValidate_ValidMaxDepth(t *testing.T) {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   0,
			SubRows: []models.Row{
				{
					Amount:     2,
					Multiplier: "x",
					Distance:   0,
					SubRows: []models.Row{
						{
							Amount:     2,
							Multiplier: "x",
							Distance:   0,
							SubRows: []models.Row{
								{
									Amount:     2,
									Multiplier: "x",
									Distance:   0,
									SubRows: []models.Row{
										{Amount: 1, Distance: 100, Content: "Test", Intensity: "GA1"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err := table.Validate()
	assert.NoError(t, err, "Valid table at max depth (4) should pass validation")
}

func TestFlattenTable(t *testing.T) {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   0,
			Content:    "Set",
			Sum:        1200,
			SubRows: []models.Row{
				{Amount: 1, Distance: 400, Content: "Kraul", Sum: 400},
				{Amount: 1, Distance: 200, Content: "Brust", Sum: 200},
			},
		},
	}
	lines := table.FlattenTable("")
	assert.Equal(t, 3, len(lines), "Should have 3 lines (1 parent + 2 child lines)")
}

func TestGetTotalVolume(t *testing.T) {
	table := models.Table{
		{Amount: 4, Distance: 100, Content: "Warmup", Sum: 400},
		{
			Amount:     8,
			Multiplier: "x",
			Distance:   0,
			Content:    "Main Set",
			Sum:        0,
			SubRows: []models.Row{
				{Amount: 1, Distance: 800, Content: "Freestyle", Sum: 800},
				{Amount: 1, Distance: 200, Content: "Kick", Sum: 200},
			},
		},
	}
	table.UpdateSum()

	total := table.GetTotalVolume()
	expected := 400 + 8000
	assert.Equal(t, expected, total, "Total volume should be correct")
}

func TestRowString_Nested(t *testing.T) {
	row := models.Row{
		Amount:     8,
		Multiplier: "x",
		Distance:   1000,
		Break:      "20",
		Content:    "Main Set",
		Intensity:  "GA1",
		Sum:        8000,
		SubRows: []models.Row{
			{Amount: 1, Distance: 800, Content: "Freestyle", Sum: 800},
			{Amount: 1, Distance: 200, Content: "Kick", Sum: 200},
		},
	}
	str := row.String()
	assert.Contains(t, str, "Set: 8 x (1 x 800m Freestyle + 1 x 200m Kick) - Main Set", "Row string should contain explicit set summary")
	assert.Contains(t, str, `↳ subrow 1/2 of "Main Set": 1 x 800m Freestyle`, "Row string should label first subrow")
	assert.Contains(t, str, `↳ subrow 2/2 of "Main Set": 1 x 200m Kick`, "Row string should label second subrow")
}

func TestRow_WithEquipment(t *testing.T) {
	row := models.Row{
		Amount:     4,
		Multiplier: "x",
		Distance:   100,
		Break:      "20",
		Content:    "Kraul-Beine",
		Intensity:  "GA1",
		Equipment:  []models.EquipmentType{models.EquipmentFins},
	}

	str := row.String()
	assert.Contains(t, str, "Flossen", "Row string should contain equipment")
}

func TestRow_EquipmentEmpty(t *testing.T) {
	row := models.Row{
		Amount:     4,
		Multiplier: "x",
		Distance:   100,
		Break:      "20",
		Content:    "Kraul",
		Intensity:  "GA1",
	}

	str := row.String()
	assert.Equal(t, "| 4 | x | 100 | 20 | Kraul | GA1 | 0 |  |", str, "Row string should have empty equipment column")
}

func TestRow_EquipmentMultiple(t *testing.T) {
	row := models.Row{
		Amount:     4,
		Multiplier: "x",
		Distance:   100,
		Break:      "20",
		Content:    "Technikübung",
		Intensity:  "TÜ",
		Equipment:  []models.EquipmentType{models.EquipmentFins, models.EquipmentBuoy},
	}

	str := row.String()
	assert.Contains(t, str, "Flossen")
	assert.Contains(t, str, "Pull buoy")
}

func TestRow_EquipmentWithSubRows(t *testing.T) {
	row := models.Row{
		Amount:     4,
		Multiplier: "x",
		Distance:   0,
		Break:      "20",
		Content:    "Main Set",
		Intensity:  "GA1",
		Sum:        400,
		Equipment:  []models.EquipmentType{models.EquipmentFins},
		SubRows: []models.Row{
			{Amount: 1, Distance: 300, Content: "Kraul", Intensity: "GA1", Sum: 300},
			{Amount: 1, Distance: 100, Content: "Kick", Intensity: "GA1", Sum: 100},
		},
	}

	str := row.String()
	assert.Contains(t, str, "Set: 4 x (1 x 300m Kraul + 1 x 100m Kick) - Main Set")
	assert.Contains(t, str, `↳ subrow 1/2 of "Main Set": 1 x 300m Kraul | break: - | intensity: GA1 | volume: 300m | equipment: -`)
	assert.Contains(t, str, "Flossen")
}

func TestScrapedPlan_MapIncludesURL(t *testing.T) {
	scrapedPlan := &models.ScrapedPlan{
		PlanID:      "test-id",
		URL:         "https://example.com/test",
		Title:       "Test Title",
		Description: "Test Description",
		Table:       models.Table{},
	}

	planMap := scrapedPlan.Map()

	assert.Contains(t, planMap, "url", "Map() should include URL")
	assert.Equal(t, "https://example.com/test", planMap["url"], "URL value should match")
	assert.Equal(t, "test-id", planMap["plan_id"], "plan_id should be present")
	assert.Equal(t, "Test Title", planMap["title"], "title should be present")
	assert.Equal(t, "Test Description", planMap["description"], "description should be present")
}
