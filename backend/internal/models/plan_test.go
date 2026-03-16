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
			Sum:        0, // Initial sum is incorrect
		},
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   50,
			Break:      "20",
			Content:    "Unterwasser-Sculling",
			Intensity:  "TÜ",
			Sum:        0, // Initial sum is incorrect
		},
		{
			Amount:     0,
			Multiplier: "",
			Distance:   0,
			Break:      "",
			Content:    "Gesamt",
			Intensity:  "",
			Sum:        0, // Initial sum is incorrect
		},
	}

	// Call UpdateSum to recalculate the sums
	table.UpdateSum()

	// Assertions
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
			Children: []models.Row{
				{Amount: 1, Distance: 800, Break: "10", Content: "Freestyle", Intensity: "GA1", Sum: 800},
				{Amount: 1, Distance: 200, Break: "0", Content: "Kick", Intensity: "GA1", Sum: 200},
			},
		},
	}
	table.UpdateSum()

	assert.Equal(t, 1000, table[0].Distance, "Parent Distance should be 1000")
	assert.Equal(t, 8000, table[0].Sum, "Parent Sum should be 8000")
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
			Children: []models.Row{
				{Amount: 1, Distance: 400, Break: "10", Content: "Kraul", Intensity: "GA2", Sum: 400},
				{Amount: 1, Distance: 100, Break: "5", Content: "Brust", Intensity: "GA2", Sum: 100},
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
			Children: []models.Row{
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

func TestValidate_ZeroAmountWithChildren(t *testing.T) {
	table := models.Table{
		{
			Amount:   0,
			Distance: 0,
			Children: []models.Row{
				{Amount: 1, Distance: 100, Content: "Test", Intensity: "GA1"},
			},
		},
	}
	err := table.Validate()
	assert.Error(t, err, "Table with zero amount and children should fail validation")
}

func TestValidate_MaxDepth(t *testing.T) {
	table := models.Table{
		{
			Amount:     2,
			Multiplier: "x",
			Distance:   0,
			Children: []models.Row{
				{
					Amount:     2,
					Multiplier: "x",
					Distance:   0,
					Children: []models.Row{
						{
							Amount:     2,
							Multiplier: "x",
							Distance:   0,
							Children: []models.Row{
								{
									Amount:     2,
									Multiplier: "x",
									Distance:   0,
									Children: []models.Row{
										{
											Amount:     2,
											Multiplier: "x",
											Distance:   0,
											Children: []models.Row{
												{
													Amount:     2,
													Multiplier: "x",
													Distance:   0,
													Children: []models.Row{
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
			Children: []models.Row{
				{
					Amount:     2,
					Multiplier: "x",
					Distance:   0,
					Children: []models.Row{
						{
							Amount:     2,
							Multiplier: "x",
							Distance:   0,
							Children: []models.Row{
								{
									Amount:     2,
									Multiplier: "x",
									Distance:   0,
									Children: []models.Row{
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
			Children: []models.Row{
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
			Children: []models.Row{
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
		Children: []models.Row{
			{Amount: 1, Distance: 800, Content: "Freestyle", Sum: 800},
			{Amount: 1, Distance: 200, Content: "Kick", Sum: 200},
		},
	}
	str := row.String()
	assert.Contains(t, str, "8x", "Row string should contain '8x'")
	assert.Contains(t, str, "800m + 200m", "Row string should contain children distances")
}
