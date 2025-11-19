package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
)

const ProfilesTableName string = "profiles"

// Retrieves a user from the database by their ID
func (db *RAGDB) GetUserProfile(ctx context.Context, id string) (*models.UserProfile, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the user with the given ID
	var user models.UserProfile
	err := pgxscan.Get(ctx, db.Conn, &user, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", ProfilesTableName), id)
	if err != nil {
		logger.Error("Error querying user", httplog.ErrAttr(err))
		return nil, fmt.Errorf("pgxscan.Select: %w", err)
	}
	return &user, nil
}

// Deletes a user and all ther associated data from the database
func (db *RAGDB) DeleteUser(ctx context.Context, id string) error {
	logger := httplog.LogEntry(ctx)

	// Delete the user from the database
	_, err := db.Conn.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		logger.Error("Error deleting user", httplog.ErrAttr(err))
		return fmt.Errorf("pgxscan.Select: %w", err)
	}
	logger.Info("User deleted successfully", "user_id", id)
	return nil
}

func (db *RAGDB) IncrementExportCount(ctx context.Context, userID, planID string) error {
	logger := httplog.LogEntry(ctx)

	if planID == "" {
		return fmt.Errorf("planID cannot be empty")
	}

	// Create a transaction
	ts, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", httplog.ErrAttr(err))
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer ts.Rollback(ctx)

	// Update the export count for the user
	if userID != "" {
		if _, err := ts.Exec(ctx,
			fmt.Sprintf(`UPDATE %s SET exports = exports + 1 WHERE user_id = $1`, ProfilesTableName),
			userID); err != nil {
			logger.Error("Error incrementing export count", httplog.ErrAttr(err))
			return fmt.Errorf("error incrementing export count: %w", err)
		}
	}

	// Update the export count for the plan
	if _, err := ts.Exec(ctx,
		fmt.Sprintf(`UPDATE %s SET exports = exports + 1 WHERE plan_id = $1`, PlanTableName),
		planID); err != nil {
		logger.Error("Error incrementing plan export count", httplog.ErrAttr(err))
		return fmt.Errorf("error incrementing plan export count: %w", err)
	}
	if err = ts.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", httplog.ErrAttr(err))
		return fmt.Errorf("error committing transaction: %w", err)
	}

	logger.Info("Export count incremented successfully", "plan_id", planID)
	return nil
}
