package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
)

type SourceOption string

const (
	// SourceOptionPlan indicates that the source is a plan.
	SourceOptionPlan SourceOption = "plan"
	// SourceOptionScraped indicates that the source is a scraped plan.
	SourceOptionScraped SourceOption = "scraped"
	// SourceOptionDonated indicates that the source is a donated plan.
	SourceOptionDonated SourceOption = "donated"
)

// Query searches for documents in the database based on the provided query and filter.
func (db *RAGDB) Query(ctx context.Context, query string, lang models.Language, filter map[string]any, method string, poolLength any) (*models.Plan, error) {
	logger := httplog.LogEntry(ctx)
	// Set the embedder to query mode
	db.Client.QueryMode()
	var docs []schema.Document
	var err error
	switch {
	case query == "" && filter == nil:
		return nil, fmt.Errorf("either a query or a filter must be provided")
	case query == "" && filter != nil:
		docs, err = db.Store.Search(ctx, 5, vectorstores.WithFilters(filter))
	case query != "" && filter == nil:
		docs, err = db.Store.SimilaritySearch(ctx, query, 5)
	case query != "" && filter != nil:
		docs, err = db.Store.SimilaritySearch(ctx, query, 5, vectorstores.WithFilters(filter))
	}
	if err != nil {
		logger.Error("Error searching for documents", httplog.ErrAttr(err))
		return nil, fmt.Errorf("error searching for documents: %w", err)
	}
	logger.Info("Documents found", "count", len(docs))
	logger.Debug("Documents:", "docs", docs)
	var plan models.Planable
	switch method {
	case "generate":
		plan, err = db.Client.GeneratePlan(ctx, query, string(lang), poolLength, docs)
		if err != nil {
			logger.Error("Error generating plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error generating plan: %w", err)
		}
	case "choose":
		if len(docs) == 0 {
			return nil, fmt.Errorf("no documents in database matching query and filters")
		}
		var planID string
		planID, err = db.Client.ChoosePlan(ctx, query, string(lang), poolLength, docs)
		if err != nil {
			logger.Error("Error choosing plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error choosing plan: %w", err)
		}

		plan, err = db.GetPlan(ctx, planID, SourceOptionPlan)
		if err != nil {
			logger.Error("Error getting plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error getting plan: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported method: %s", method)
	}

	genericPlan := plan.Plan()

	if lang != "de" {
		genericPlan, err = db.Client.TranslatePlan(ctx, genericPlan, lang)
		if err != nil {
			logger.Error("Error translating plan", httplog.ErrAttr(err))
			return nil, fmt.Errorf("error translating plan: %w", err)
		}

	}

	logger.Debug("Plan generated successfully", "plan", genericPlan)

	return genericPlan, nil
}

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

func (db *RAGDB) UpsertPlan(ctx context.Context, plan models.Plan, userID string) error {
	logger := httplog.LogEntry(ctx)

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
		return fmt.Errorf("failed to check plan existence: %w", err)
	}

	// If it doesn't exist, create a new plan id
	if !exists {
		plan.PlanID = GenerateRandomUUID()
	}

	// Start a transaction
	tx, err := db.Conn.Begin(ctx)
	if err != nil {
		logger.Error("Error starting transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Upsert the plan and add it to the user's history
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
		return fmt.Errorf("failed to upsert plan: %w", err)
	}

	// Add the plan to the user's history
	_, err = tx.Exec(ctx,
		`INSERT INTO history (user_id, plan_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, plan.PlanID,
	)
	if err != nil {
		logger.Error("Error adding plan to user history", httplog.ErrAttr(err))
		return fmt.Errorf("failed to add plan to user history: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		logger.Error("Error committing transaction", httplog.ErrAttr(err))
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	logger.Info("Plan upserted successfully", "plan_id", plan.PlanID)
	return nil
}
