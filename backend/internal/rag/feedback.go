package rag

import (
	"context"
	"errors"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
)

const FeedbackTable string = "feedback"

var ErrFeedbackPlanNotFound = errors.New("feedback plan not found")

const feedbackColumns = `user_id, plan_id, rating, comment, created_at, updated_at, was_swam, difficulty_rating, coalesce(removed_from_history, false) as removed_from_history`

// SubmitFeedback atomically authorizes and upserts feedback for a plan.
func (db *RAGDB) SubmitFeedback(ctx context.Context, feedback *models.Feedback) error {
	logger := httplog.LogEntry(ctx)
	var difficultyRating any
	if feedback.DifficultyRating != nil {
		difficultyRating = *feedback.DifficultyRating
	}

	var submitted bool
	err := db.Conn.QueryRow(ctx, `
		select public.submit_feedback($1, $2, $3, $4, $5, $6)
	`, feedback.UserID, feedback.PlanID, feedback.Rating, feedback.WasSwam, difficultyRating, feedback.Comment).Scan(&submitted)
	if err != nil {
		logger.Error("Error submitting feedback", httplog.ErrAttr(err))
		return err
	}
	if !submitted {
		return ErrFeedbackPlanNotFound
	}

	logger.Debug("Feedback submitted successfully", "plan_id", feedback.PlanID)
	return nil
}

// Get all feedback for a plan from the database
func (db *RAGDB) GetAllFeedbackForPlan(ctx context.Context, planID string) ([]*models.Feedback, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for all feedback for the given plan ID
	var feedbacks []*models.Feedback
	err := pgxscan.Select(ctx, db.Conn, &feedbacks, fmt.Sprintf("SELECT %s FROM %s WHERE plan_id = $1", feedbackColumns, FeedbackTable), planID)
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
	err := pgxscan.Select(ctx, db.Conn, &feedbacks, fmt.Sprintf("SELECT %s FROM %s WHERE user_id = $1", feedbackColumns, FeedbackTable), userID)
	if err != nil {
		logger.Error("Error querying feedback", httplog.ErrAttr(err))
		return nil, err
	}
	return feedbacks, nil
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
	logger.Debug("Feedback deleted successfully", "user_id", userID, "plan_id", planID)
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
	logger.Debug("All feedback deleted successfully", "user_id", userID)
	return nil
}
