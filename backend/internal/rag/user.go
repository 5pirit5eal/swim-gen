package rag

import (
	"context"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/go-chi/httplog/v2"
)

// Adds a new user to the database
func (db *RAGDB) AddUser(ctx context.Context, user *models.User) error {
	logger := httplog.LogEntry(ctx)

	// Create a new user in the database using the struct fields
	_, err := db.Conn.Exec(ctx, "INSERT INTO users (user_id, name, email, created_at, last_active) VALUES ($1, $2, $3, $4, $5)",
		GenerateUUID(user.Email, user.Name), user.Name, user.Email, user.CreatedAt, user.LastActive)
	if err != nil {
		logger.Error("Error creating user", httplog.ErrAttr(err))
		return err
	}
	logger.Info("User created successfully", "user", user)
	return nil
}

// Retrieves a user from the database by their ID
func (db *RAGDB) GetUser(ctx context.Context, id string) (*models.User, error) {
	logger := httplog.LogEntry(ctx)

	// Query the database for the user with the given ID
	var user models.User
	err := pgxscan.Get(ctx, db.Conn, &user, "SELECT * FROM users WHERE id = $1", id)
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
