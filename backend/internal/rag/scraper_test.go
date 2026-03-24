package rag_test

import (
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestToLangChainDoc_ValidPlan(t *testing.T) {
	table := models.Table{
		{
			Amount:     4,
			Multiplier: "x",
			Distance:   100,
			Break:      "20",
			Content:    "Freestyle",
			Intensity:  "GA1",
			Sum:        400,
		},
	}

	plan := models.ScrapedPlan{
		PlanID:      "test-plan-id",
		URL:         "https://example.com/test",
		Title:       "Test Plan",
		Description: "Test Description",
		Table:       table,
	}

	doc := models.Document{
		Plan: &plan,
		Meta: &models.Metadata{
			Reasoning: "Test reasoning",
		},
	}

	result, err := doc.ToLangChainDoc()
	assert.NoError(t, err)

	assert.NotEmpty(t, result.PageContent)
	assert.Equal(t, "test-plan-id", result.Metadata["plan_id"])
	assert.Equal(t, "https://example.com/test", result.Metadata["url"])
	assert.Equal(t, "Test Plan", result.Metadata["title"])
}

func TestToLangChainDoc_EmptyPlanID(t *testing.T) {
	plan := models.ScrapedPlan{
		URL:         "https://example.com/test",
		Title:       "Test Plan",
		Description: "Test Description",
		Table:       models.Table{},
	}

	doc := models.Document{
		Plan: &plan,
		Meta: &models.Metadata{},
	}

	result, err := doc.ToLangChainDoc()
	assert.NoError(t, err, "Should not return error when plan_id key exists (even if empty)")
	assert.Equal(t, "", result.Metadata["plan_id"], "plan_id should be empty string")
}
