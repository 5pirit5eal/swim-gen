package rag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementExportCountValidation(t *testing.T) {
	db := &RAGDB{}
	ctx := context.Background()

	t.Run("Empty planID returns nil without DB interaction", func(t *testing.T) {
		err := db.IncrementExportCount(ctx, "00000000-0000-0000-0000-000000000001", "")
		assert.NoError(t, err)
	})

	t.Run("Invalid planID returns error", func(t *testing.T) {
		err := db.IncrementExportCount(ctx, "", "invalid-plan-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid plan UUID")
	})

	t.Run("Invalid userID returns error", func(t *testing.T) {
		err := db.IncrementExportCount(ctx, "invalid-user-uuid", "00000000-0000-0000-0000-000000000001")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user UUID")
	})
}
