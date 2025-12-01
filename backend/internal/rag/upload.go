package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
)

const DonatedPlanTable string = "donations"

// Add uploaded plan to the database
func (db *RAGDB) AddUploadedPlan(ctx context.Context, upload *models.DonatedPlan) error {
	logger := httplog.LogEntry(ctx)

	// Begin transaction for plan and uploaded entry
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error beginning transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Add the plan to the plans table
	_, err = tx.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (plan_id, title, description, plan_table)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (plan_id) DO NOTHING`, PlanTableName),
		upload.PlanID, upload.Title, upload.Description, upload.Table)
	if err != nil {
		logger.Error("Error inserting plan", httplog.ErrAttr(err))
		return fmt.Errorf("failed to insert plan: %w", err)
	}

	// Create a new donation entry in the database using the struct fields
	_, err = tx.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (user_id, plan_id, allow_sharing) VALUES ($1, $2, $3)", DonatedPlanTable),
		upload.UserID, upload.PlanID, upload.AllowSharing)
	if err != nil {
		logger.Error("Error creating donation", httplog.ErrAttr(err))
		return fmt.Errorf("failed to insert donation: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Debug("Uploaded plan added successfully", "plan", upload)
	return nil
}

// Get uploaded plans for a user
func (db *RAGDB) GetUploadedPlans(ctx context.Context, userID string) ([]*models.DonatedPlan, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the uploaded plans
	var plans []*models.DonatedPlan
	err := pgxscan.Select(ctx, db.Conn, &plans,
		fmt.Sprintf(`
			SELECT dp.user_id, dp.plan_id, dp.created_at, dp.allow_sharing, p.title, p.description, p.plan_table
			FROM %s dp
			JOIN %s p ON dp.plan_id = p.plan_id
			WHERE dp.user_id = $1
		`, DonatedPlanTable, PlanTableName), userID)
	if err != nil {
		logger.Error("Error querying uploaded plans", httplog.ErrAttr(err))
		return nil, err
	}

	if len(plans) == 0 {
		plans = []*models.DonatedPlan{}
	}

	logger.Debug("Uploaded plans retrieved successfully", "count", len(plans))
	return plans, nil
}

// Get a single uploaded plan by plan ID
func (db *RAGDB) GetUploadedPlan(ctx context.Context, planID string) (*models.DonatedPlan, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the uploaded plan
	var plan models.DonatedPlan
	err := pgxscan.Get(ctx, db.Conn, &plan,
		fmt.Sprintf(`
			SELECT dp.user_id, dp.plan_id, dp.created_at, dp.allow_sharing, p.title, p.description, p.plan_table
			FROM %s dp
			JOIN %s p ON dp.plan_id = p.plan_id
			WHERE dp.plan_id = $1`, DonatedPlanTable, PlanTableName), planID)
	if err != nil {
		logger.Error("Error querying uploaded plan", httplog.ErrAttr(err))
		return nil, err
	}

	logger.Debug("Uploaded plan retrieved successfully", "plan", plan)
	return &plan, nil
}
