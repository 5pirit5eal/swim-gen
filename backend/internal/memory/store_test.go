package memory

import (
	"context"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test requires a running database connection, so we skip it if no connection string is provided
// or if we are in a CI/CD environment without a DB.
// For now, this serves as a compilation check and a template for integration testing.
func TestMemoryStore_AddMessage(t *testing.T) {
	// Skip if no DB connection available (mock check)
	// if os.Getenv("TEST_DB_URL") == "" {
	// 	t.Skip("Skipping test; TEST_DB_URL not set")
	// }
	t.Skip("Skipping integration test requiring DB")

	ctx := context.Background()
	// db, err := pgxpool.New(ctx, os.Getenv("TEST_DB_URL"))
	// require.NoError(t, err)
	// defer db.Close()
	var db *pgxpool.Pool // Mock or nil for now

	store := NewMemoryStore(db)

	planID := "test-plan-id"
	userID := "test-user-id"
	role := models.RoleUser
	content := "Hello AI"

	msg, err := store.AddMessage(ctx, planID, userID, role, content, nil, nil)
	if db == nil {
		// Expected to fail with nil pointer if we actually ran it
		return
	}
	require.NoError(t, err)
	assert.NotEmpty(t, msg.ID)
	assert.Equal(t, content, msg.Content)

	// Test DeleteConversation (Compilation check)
	err = store.DeleteConversation(ctx, planID)
	require.NoError(t, err)

	// Test DeleteMessage (Compilation check)
	err = store.DeleteMessage(ctx, msg.ID)
	require.NoError(t, err)

	// Test UpdateMessage (Compilation check)
	err = store.UpdateMessage(ctx, msg.ID, "New Content", nil)
	require.NoError(t, err)

	// Test DeleteMessagesAfter (Compilation check)
	err = store.DeleteMessagesAfter(ctx, msg.ID)
	require.NoError(t, err)
}
