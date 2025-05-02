package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
)

// Add donated plan to the database
func (db *RAGDB) AddDonatedPlan(ctx context.Context, donation *models.DonatedPlan, meta *models.Metadata) error {
	logger := httplog.LogEntry(ctx)

	// Add the donated plan to the vector store
	doc, err := models.PlanToDoc(&models.Document{Plan: donation, Meta: meta})
	if err != nil {
		logger.Error("Error converting plan to document", httplog.ErrAttr(err))
		return fmt.Errorf("PlanToDoc: %w", err)
	}
	_, err = db.Store.AddDocuments(ctx, []schema.Document{doc})
	if err != nil {
		logger.Error("Error adding donation to vector store", httplog.ErrAttr(err))
		return fmt.Errorf("Store.AddDocuments: %w", err)
	}

	// Create a new donation entry in the database using the struct fields
	_, err = db.Conn.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (user_id, plan_id, created_at) VALUES ($1, $2, $3)", DonatedPlanTable),
		donation.UserID, donation.PlanID, donation.CreatedAt)
	if err != nil {
		logger.Error("Error creating donation", httplog.ErrAttr(err))
		return err
	}

	logger.Info("Donation added successfully", "donation", donation)
	return nil
}

// Get donated plans for a user
func (db *RAGDB) GetDonatedPlans(ctx context.Context, userID string) ([]*models.DonatedPlan, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the donated plans
	var plans []*models.DonatedPlan
	err := pgxscan.Select(ctx, db.Conn, &plans,
		fmt.Sprintf(`
			SELECT dp.user_id, dp.plan_id, dp.created_at, p.title, p.description, p.plan_table
			FROM %s dp
			JOIN %s p ON dp.plan_id = p.plan_id
			WHERE dp.user_id = $1
		`, DonatedPlanTable, PlanTableName), userID)
	if err != nil {
		logger.Error("Error querying donated plans", httplog.ErrAttr(err))
		return nil, err
	}

	logger.Info("Donated plans retrieved successfully", "count", len(plans))
	return plans, nil
}

// Get a single donated plan by plan ID
func (db *RAGDB) GetDonatedPlan(ctx context.Context, planID string) (*models.DonatedPlan, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the donated plan
	var plan models.DonatedPlan
	err := pgxscan.Get(ctx, db.Conn, &plan,
		fmt.Sprintf(`
			SELECT dp.user_id, dp.plan_id, dp.created_at, p.title, p.description, p.plan_table
			FROM %s dp
			JOIN %s p ON dp.plan_id = p.plan_id
			WHERE dp.plan_id = $1`, DonatedPlanTable, PlanTableName), planID)
	if err != nil {
		logger.Error("Error querying donated plan", httplog.ErrAttr(err))
		return nil, err
	}

	logger.Info("Donated plan retrieved successfully", "plan", plan)
	return &plan, nil
}
