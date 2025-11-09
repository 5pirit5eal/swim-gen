package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
)

const HistoryTableName string = "history"

func (db *RAGDB) GetPlan(ctx context.Context, planID string, source SourceOption) (models.Planable, error) {
	logger := httplog.LogEntry(ctx)
	// Query the database for the plan with the given ID
	switch source {
	case SourceOptionPlan:
		var plan models.Plan
		err := pgxscan.Get(ctx, db.Conn, &plan, fmt.Sprintf(`SELECT plan_id, title, description, plan_table FROM %s WHERE plan_id = $1`, PlanTableName), planID)
		if err != nil {
			logger.Error("Error querying plan", httplog.ErrAttr(err))
			return nil, err
		}
		return &plan, nil
	case SourceOptionScraped:
		return db.GetScrapedPlan(ctx, planID)
	case SourceOptionDonated:
		return db.GetDonatedPlan(ctx, planID)
	}
	return nil, fmt.Errorf("unsupported source option: %s", source)
}

func (db *RAGDB) UpsertPlan(ctx context.Context, plan models.Plan, userID string) (string, error) {
	logger := httplog.LogEntry(ctx)

	if plan.PlanID == "" {
		logger.Info("No plan ID provided, generating a new one.")
		plan.PlanID = uuid.New().String()
	}

	// Check if the plan exists and the user owns it by querying the donation and history table
	var exists bool
	err := pgxscan.Get(ctx, db.Conn, &exists, fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1
            FROM %s h
            JOIN %s d ON h.plan_id = d.plan_id
            WHERE h.plan_id = $1 AND h.user_id = $2 AND d.user_id = $2
        )`, HistoryTableName, DonatedPlanTable), plan.PlanID, userID)

	if err != nil {
		logger.Error("Error checking plan existence", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to check plan existence: %w", err)
	}

	// If it doesn't exist, create a new plan id
	if !exists {
		plan.PlanID = uuid.New().String()
	}

	// Start a transaction to upsert the plan and add it to the user's history
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Add the plan to the plans table
	_, err = tx.Exec(ctx,
		fmt.Sprintf(`
            INSERT INTO %s (plan_id, title, description, plan_table)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (plan_id) DO UPDATE
            SET title = EXCLUDED.title,
                description = EXCLUDED.description,
                plan_table = EXCLUDED.plan_table,
                updated_at = now()
        `, PlanTableName),
		plan.PlanID, plan.Title, plan.Description, plan.Table,
	)
	if err != nil {
		logger.Error("Error upserting plan", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to upsert plan: %w", err)
	}

	// Add the plan to the user's history
	_, err = tx.Exec(ctx,
		`INSERT INTO history (user_id, plan_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, plan.PlanID,
	)
	if err != nil {
		logger.Error("Error adding plan to user history", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to add plan to user history: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Plan upserted successfully", "plan_id", plan.PlanID)
	return plan.PlanID, nil
}

func (db *RAGDB) AddPlanToHistory(ctx context.Context, plan *models.Plan, userID string) error {
	logger := httplog.LogEntry(ctx)

	// Check that the user exists
	userProfile, err := db.GetUserProfile(ctx, userID)
	if err != nil {
		logger.Error("Error retrieving user", httplog.ErrAttr(err))
		return fmt.Errorf("error retrieving user: %w", err)
	}

	// Start a transaction
	ts, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", httplog.ErrAttr(err))
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer ts.Rollback(ctx)

	// Insert the plan into the plans table
	if _, err := ts.Exec(ctx, fmt.Sprintf(`
        INSERT INTO %s (plan_id, title, description, plan_table)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (plan_id) DO NOTHING`, PlanTableName),
		plan.PlanID, plan.Title, plan.Description, plan.Table); err != nil {
		logger.Error("Error inserting plan", httplog.ErrAttr(err))
		return fmt.Errorf("error inserting plan: %w", err)
	}

	// Insert the plan into the user's history
	if _, err = ts.Exec(ctx, fmt.Sprintf(
		`INSERT INTO %s (user_id, plan_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, HistoryTableName),
		userProfile.UserID, plan.PlanID); err != nil {
		logger.Error("Error adding plan to user history", httplog.ErrAttr(err))
		return fmt.Errorf("error adding plan to user history: %w", err)
	}

	if err = ts.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", httplog.ErrAttr(err))
		return fmt.Errorf("error committing transaction: %w", err)
	}

	logger.Info("Plan added to user history successfully", "user_id", userID, "plan_id", plan.PlanID)
	return nil
}
