package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
)

// Add donated plan to the database
func (db *RAGDB) AddDonatedPlan(ctx context.Context, donation *models.DonatedPlan, meta *models.Metadata) error {
	logger := httplog.LogEntry(ctx)

	// Add the donated plan to the vector store
	_, err := db.Store.AddDocuments(ctx, []schema.Document{models.PlanToDoc(donation, meta)})
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
