package models_test

import (
	"encoding/json"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestMetadataSchema(t *testing.T) {
	schema, err := models.MetadataSchema()
	assert.NoError(t, err, "Failed to retrieve schema")

	// check if the schema is valid json
	var result map[string]interface{}
	err = json.Unmarshal([]byte(schema), &result)
	assert.NoError(t, err, "Failed to unmarshal schema")
	assert.NotEmpty(t, result, "Schema should not be empty")
}

func TestEquipmentConstants(t *testing.T) {
	// Test that all constants are properly defined with German values
	assert.Equal(t, "Flossen", string(models.EquipmentFins))
	assert.Equal(t, "Kickboard", string(models.EquipmentKickboard))
	assert.Equal(t, "Handpaddles", string(models.EquipmentPaddles))
	assert.Equal(t, "Pull buoy", string(models.EquipmentBuoy))
	assert.Equal(t, "Schnorchel", string(models.EquipmentSnorkel))
}
