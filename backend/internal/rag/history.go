package rag

import (
	"context"
	"fmt"

	"github.com/go-chi/httplog/v2"
)

const HistoryTableName string = "history"

func (db *RAGDB) AddPlanToHistory(ctx context.Context, userID, planID string) error {
	logger := httplog.LogEntry(ctx)

	// Check that the user exists
	userProfile, err := db.GetUserProfile(ctx, userID)
	if err != nil {
		logger.Error("Error retrieving user", httplog.ErrAttr(err))
		return fmt.Errorf("error retrieving user: %w", err)
	}

	// Insert the plan into the user's history
	_, err = db.Conn.Exec(ctx,
		fmt.Sprintf(`INSERT INTO %s (user_id, plan_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, HistoryTableName),
		userProfile.UserID, planID)
	if err != nil {
		logger.Error("Error adding plan to user history", httplog.ErrAttr(err))
		return fmt.Errorf("error adding plan to user history: %w", err)
	}

	logger.Info("Plan added to user history successfully", "user_id", userID, "plan_id", planID)
	return nil
}
