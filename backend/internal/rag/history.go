package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
)

const (
	HistoryTableName string = "history"
	PlanTableName    string = "plans"
)

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
		return db.GetUploadedPlan(ctx, planID)
	}
	return nil, fmt.Errorf("unsupported source option: %s", source)
}

func (db *RAGDB) UpsertPlan(ctx context.Context, plan models.Plan, userID string) (string, error) {
	logger := httplog.LogEntry(ctx)

	if plan.PlanID == "" {
		logger.Debug("No plan ID provided, generating a new one.")
		plan.PlanID = uuid.New().String()
	}

	// Check if the plan exists and the user owns it by querying the donation and history table
	// The plan is owned by the user if it is in either their donated plans or their history
	var exists bool
	err := pgxscan.Get(ctx, db.Conn, &exists, fmt.Sprintf(`
        SELECT EXISTS (
            SELECT 1
            FROM (
                SELECT plan_id, user_id FROM %s
                UNION ALL
                SELECT plan_id, user_id FROM %s
            ) as combined_plans
            WHERE plan_id = $1 AND user_id = $2
        )`, HistoryTableName, DonatedPlanTable), plan.PlanID, userID)

	if err != nil {
		logger.Error("Error checking plan existence", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to check plan existence: %w", err)
	}

	// If it doesn't exist, create a new plan id
	if !exists {
		logger.Debug("Plan does not exist for user, generating new plan ID")
		plan.PlanID = uuid.New().String()
	}

	// Start a transaction to upsert the plan and add it to the user's history
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", httplog.ErrAttr(err))
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Add the plan to the plans table
	logger.Debug("Upserting plan into plans table")
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
	logger.Debug("Adding plan to user history")
	if !exists {
		_, err = tx.Exec(ctx,
			`INSERT INTO history (user_id, plan_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			userID, plan.PlanID,
		)
		if err != nil {
			logger.Error("Error adding plan to user history", httplog.ErrAttr(err))
			return "", fmt.Errorf("failed to add plan to user history: %w", err)
		}
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
	defer func() { _ = ts.Rollback(ctx) }()

	// Insert the plan into the plans table
	if _, err := ts.Exec(ctx, fmt.Sprintf(`
        INSERT INTO %s (plan_id, title, description, plan_table)
		VALUES ($1, $2, $3, $4)`, PlanTableName),
		plan.PlanID, plan.Title, plan.Description, plan.Table); err != nil {
		logger.Error("Error inserting plan", httplog.ErrAttr(err))
		return fmt.Errorf("error inserting plan: %w", err)
	}

	// Insert the plan into the user's history
	if _, err = ts.Exec(ctx, fmt.Sprintf(
		`INSERT INTO %s (user_id, plan_id) VALUES ($1, $2)`, HistoryTableName),
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

func (db *RAGDB) SharePlan(ctx context.Context, planID, userID string, method models.SharingMethod) (string, error) {
	logger := httplog.LogEntry(ctx)

	// Calculate a short uuid for the shared plan based on planID and userID
	urlHash := uuid.NewSHA1(uuid.NameSpaceURL, []byte(planID+userID)).String()

	switch method {
	case models.SharingMethodLink:
		// Insert the shared plan into the shared_plans table
		row := db.Conn.QueryRow(ctx,
			fmt.Sprintf(`
                INSERT INTO %s (user_id, plan_id, url_hash)
                VALUES ($1, $2, $3)
                ON CONFLICT (plan_id) DO UPDATE
                SET url_hash = %s.url_hash
                RETURNING url_hash
            `, "shared_plans", "shared_plans"),
			userID, planID, urlHash,
		)
		err := row.Scan(&urlHash)
		if err != nil {
			logger.Error("Error sharing plan", httplog.ErrAttr(err))
			return "", fmt.Errorf("failed to share plan: %w", err)
		}

		logger.Info("Plan shared successfully", "plan_id", planID, "user_id", userID, "url_hash", urlHash)
		return urlHash, nil
	case models.SharingMethodEmail:
		// proceed
		return "", fmt.Errorf("email sharing not implemented yet")
	default:
		return "", fmt.Errorf("unsupported sharing method: %s", method)
	}
}

// DeletePlan removes a plan from the user's history.
// If the plan has feedback, it preserves the plan data and marks the feedback as removed_from_history.
// If no feedback exists, it deletes the plan entirely (CASCADE removes related data).
func (db *RAGDB) DeletePlan(ctx context.Context, planID, userID string) error {
	logger := httplog.LogEntry(ctx)

	// Start a transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", httplog.ErrAttr(err))
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Check if the plan exists and user owns it (via history, donations, or shared_plans)
	var exists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM history WHERE plan_id = $1 AND user_id = $2
			UNION
			SELECT 1 FROM donations WHERE plan_id = $1 AND user_id = $2
			UNION
			SELECT 1 FROM shared_plans WHERE plan_id = $1 AND user_id = $2
		)`,
		planID, userID,
	).Scan(&exists)

	if err != nil {
		logger.Error("Error checking plan existence", httplog.ErrAttr(err))
		return fmt.Errorf("failed to check plan existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("plan not found in user history or user does not own the plan")
	}

	// Check if feedback exists for this plan
	var hasFeedback bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM feedback WHERE plan_id = $1)`,
		planID,
	).Scan(&hasFeedback)

	if err != nil {
		logger.Error("Error checking feedback existence", httplog.ErrAttr(err))
		return fmt.Errorf("failed to check feedback existence: %w", err)
	}

	if hasFeedback {
		// Plan has feedback - preserve the plan, only remove from history and mark feedback
		logger.Info("Plan has feedback, preserving plan data", "plan_id", planID)

		// Remove from user's history only
		_, err = tx.Exec(ctx,
			`DELETE FROM history WHERE plan_id = $1 AND user_id = $2`,
			planID, userID,
		)
		if err != nil {
			logger.Error("Error removing plan from history", httplog.ErrAttr(err))
			return fmt.Errorf("failed to remove plan from history: %w", err)
		}

		// Remove from shared_history of any other users who have this plan shared with them
		_, err = tx.Exec(ctx,
			`DELETE FROM shared_history WHERE plan_id = $1`,
			planID,
		)
		if err != nil {
			logger.Error("Error removing plan from shared_history", httplog.ErrAttr(err))
			return fmt.Errorf("failed to remove plan from shared_history: %w", err)
		}

		// Remove from shared_plans table
		_, err = tx.Exec(ctx,
			`DELETE FROM shared_plans WHERE plan_id = $1`,
			planID,
		)
		if err != nil {
			logger.Error("Error removing plan from shared_plans", httplog.ErrAttr(err))
			return fmt.Errorf("failed to remove plan from shared_plans: %w", err)
		}

		// Mark feedback as removed_from_history
		_, err = tx.Exec(ctx,
			`UPDATE feedback SET removed_from_history = true WHERE plan_id = $1 AND user_id = $2`,
			planID, userID,
		)
		if err != nil {
			logger.Error("Error marking feedback as removed", httplog.ErrAttr(err))
			return fmt.Errorf("failed to mark feedback as removed: %w", err)
		}

		logger.Info("Plan removed from history, feedback preserved", "plan_id", planID, "user_id", userID)
	} else {
		// No feedback - delete the plan entirely
		// This will CASCADE to: history, conversation, shared_plans, shared_history
		_, err = tx.Exec(ctx,
			fmt.Sprintf(`DELETE FROM %s WHERE plan_id = $1`, PlanTableName),
			planID,
		)
		if err != nil {
			logger.Error("Error deleting plan", httplog.ErrAttr(err))
			return fmt.Errorf("failed to delete plan: %w", err)
		}

		logger.Info("Plan deleted successfully", "plan_id", planID, "user_id", userID)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
