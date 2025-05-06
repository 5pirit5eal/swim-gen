package models_test

import (
	"encoding/json"
	"testing"

	"github.com/5pirit5eal/swim-rag/internal/models"
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

func TestTableSchema(t *testing.T) {
	schema, err := models.TableSchema()
	assert.NoError(t, err, "Failed to retrieve schema")

	// check if the schema is valid json
	var result map[string]interface{}
	err = json.Unmarshal([]byte(schema), &result)
	assert.NoError(t, err, "Failed to unmarshal schema")
	assert.NotEmpty(t, result, "Schema should not be empty")
}
