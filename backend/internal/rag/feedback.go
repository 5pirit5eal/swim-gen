package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
)

const FeedbackTable string = "feedback"

// Add user feedback for an existing plan to the database
func (db *RAGDB) AddFeedback(ctx context.Context, feedback *models.Feedback) error {
	logger := httplog.LogEntry(ctx)

	// Create a new feedback entry in the database using the struct fields
	_, err := db.Conn.Exec(ctx,
		fmt.Sprintf("INSERT INTO %s (user_id, plan_id, rating, was_swam, difficulty_rating, comment) VALUES ($1, $2, $3, $4, $5, $6)", FeedbackTable),
		feedback.UserID, feedback.PlanID, feedback.Rating, feedback.WasSwam, feedback.DifficultyRating, feedback.Comment)
	if err != nil {
		logger.Error("Error creating feedback", httplog.ErrAttr(err))
		return err
	}
	logger.Info("Feedback created successfully", "feedback", feedback)
	return nil
}

// Get user feedback for a plan from the database
func (db *RAGDB) GetFeedback(ctx context.Context, userID string, planID string) (*models.Feedback, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the feedback with the given user ID and plan ID
	var feedback models.Feedback
	err := db.Conn.QueryRow(ctx, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND plan_id = $2", FeedbackTable), userID, planID).Scan(
		&feedback.UserID,
		&feedback.PlanID,
		&feedback.Rating,
		&feedback.Comment,
		&feedback.CreatedAt,
		&feedback.UpdatedAt,
		&feedback.WasSwam,
		&feedback.DifficultyRating,
	)
	if err != nil {
		logger.Error("Error querying feedback", httplog.ErrAttr(err))
		return nil, err
	}
	return &feedback, nil
}

// Get all feedback for a plan from the database
func (db *RAGDB) GetAllFeedbackForPlan(ctx context.Context, planID string) ([]*models.Feedback, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for all feedback for the given plan ID
	var feedbacks []*models.Feedback
	err := pgxscan.Select(ctx, db.Conn, &feedbacks, fmt.Sprintf("SELECT * FROM %s WHERE plan_id = $1", FeedbackTable), planID)
	if err != nil {
		logger.Error("Error querying feedback", httplog.ErrAttr(err))
		return nil, err
	}
	return feedbacks, nil
}

// Get all feedback from a user from the database
func (db *RAGDB) GetAllFeedbackFromUser(ctx context.Context, userID string) ([]*models.Feedback, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for all feedback from the given user ID
	var feedbacks []*models.Feedback
	err := pgxscan.Select(ctx, db.Conn, &feedbacks, fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", FeedbackTable), userID)
	if err != nil {
		logger.Error("Error querying feedback", httplog.ErrAttr(err))
		return nil, err
	}
	return feedbacks, nil
}

// Update user feedback for a plan in the database
func (db *RAGDB) UpdateFeedback(ctx context.Context, feedback *models.Feedback) error {
	logger := httplog.LogEntry(ctx)

	// Update the feedback entry in the database using the struct fields
	_, err := db.Conn.Exec(ctx, fmt.Sprintf("UPDATE %s SET rating = $1, was_swam = $2, difficulty_rating = $3, comment = $4 WHERE user_id = $5 AND plan_id = $6", FeedbackTable),
		feedback.Rating, feedback.WasSwam, feedback.DifficultyRating, feedback.Comment, feedback.UserID, feedback.PlanID)
	if err != nil {
		logger.Error("Error updating feedback", httplog.ErrAttr(err))
		return err
	}
	logger.Info("Feedback updated successfully", "feedback", feedback)
	return nil
}

// Delete user feedback for a plan from the database
func (db *RAGDB) DeleteFeedback(ctx context.Context, userID string, planID string) error {
	logger := httplog.LogEntry(ctx)

	// Delete the feedback entry from the database
	_, err := db.Conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND plan_id = $2", FeedbackTable), userID, planID)
	if err != nil {
		logger.Error("Error deleting feedback", httplog.ErrAttr(err))
		return err
	}
	logger.Info("Feedback deleted successfully", "user_id", userID, "plan_id", planID)
	return nil
}

// Delete all feedback by a user from the database
func (db *RAGDB) DeleteAllFeedbackFromUser(ctx context.Context, userID string) error {
	logger := httplog.LogEntry(ctx)

	// Delete all feedback entries from the database for the given user ID
	_, err := db.Conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", FeedbackTable), userID)
	if err != nil {
		logger.Error("Error deleting feedback", httplog.ErrAttr(err))
		return err
	}
	logger.Info("All feedback deleted successfully", "user_id", userID)
	return nil
}
