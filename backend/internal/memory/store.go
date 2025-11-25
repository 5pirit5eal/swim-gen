package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const MemoryTableName = "memory"

type Store interface {
	AddMessage(ctx context.Context, planID, userID string, role models.Role, content string, previousMessageID *string, planSnapshot *models.Plan) (*models.Message, error)
	GetConversation(ctx context.Context, planID string) ([]models.Message, error)
	GetLastMessage(ctx context.Context, q pgxscan.Querier, planID string) (*models.Message, error)
	DeleteConversation(ctx context.Context, planID string) error
	DeleteMessage(ctx context.Context, messageID string) error
	UpdateMessage(ctx context.Context, messageID, content string, planSnapshot *models.Plan) error
	DeleteMessagesAfter(ctx context.Context, messageID string) error
}

type MemoryStore struct {
	db *pgxpool.Pool
}

// Ensure MemoryStore implements Store
var _ Store = (*MemoryStore)(nil)

func NewMemoryStore(db *pgxpool.Pool) Store {
	return &MemoryStore{db: db}
}

// AddMessage inserts a new message into the memory table and updates the linked list.
func (s *MemoryStore) AddMessage(ctx context.Context, planID, userID string, role models.Role, content string, previousMessageID *string, planSnapshot *models.Plan) (*models.Message, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// If PreviousMessageID is not provided, try to find the last message in the conversation
	if previousMessageID == nil {
		lastMsg, err := s.GetLastMessage(ctx, tx, planID)
		if err != nil {
			return nil, fmt.Errorf("failed to get last message: %w", err)
		}
		if lastMsg != nil {
			previousMessageID = &lastMsg.ID
		}
	}

	// Insert the new message
	var newMessageID string
	query := fmt.Sprintf(`
		INSERT INTO %s (plan_id, user_id, role, content, previous_message_id, plan_snapshot)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, MemoryTableName)
	err = tx.QueryRow(ctx, query,
		planID,
		userID,
		role,
		content,
		previousMessageID,
		planSnapshot,
	).Scan(&newMessageID)

	if err != nil {
		return nil, fmt.Errorf("failed to insert message: %w", err)
	}

	// Update the previous message's next_message_id if it exists
	if previousMessageID != nil {
		updateQuery := fmt.Sprintf(`
			UPDATE %s
			SET next_message_id = $1
			WHERE id = $2
		`, MemoryTableName)
		_, err = tx.Exec(ctx, updateQuery, newMessageID, *previousMessageID)
		if err != nil {
			return nil, fmt.Errorf("failed to update previous message: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.Message{
		ID:                newMessageID,
		PlanID:            planID,
		UserID:            userID,
		Role:              role,
		Content:           content,
		PreviousMessageID: previousMessageID,
		PlanSnapshot:      planSnapshot,
	}, nil
}

// GetConversation retrieves the full conversation for a plan, ordered by the linked list.
// Note: For large conversations, we might want to paginate or limit this.
// For now, we fetch all and sort in Go or use a recursive CTE.
// Recursive CTE is better for ordering linked lists in SQL.
func (s *MemoryStore) GetConversation(ctx context.Context, planID string) ([]models.Message, error) {
	query := `
		WITH RECURSIVE conversation AS (
			SELECT id, plan_id, user_id, role, content, previous_message_id, next_message_id, plan_snapshot, created_at
			FROM memory
			WHERE plan_id = $1 AND previous_message_id IS NULL

			UNION ALL

			SELECT m.id, m.plan_id, m.user_id, m.role, m.content, m.previous_message_id, m.next_message_id, m.plan_snapshot, m.created_at
			FROM memory m
			INNER JOIN conversation c ON m.previous_message_id = c.id
		)
		SELECT * FROM conversation;
	`

	var messages []models.Message
	if err := pgxscan.Select(ctx, s.db, &messages, query, planID); err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return messages, nil
}

// GetLastMessage retrieves the last message in the conversation (where next_message_id is NULL).
// It accepts a querier (tx or pool) to support transactions.
func (s *MemoryStore) GetLastMessage(ctx context.Context, q pgxscan.Querier, planID string) (*models.Message, error) {
	query := `
		SELECT id, plan_id, user_id, role, content, previous_message_id, next_message_id, plan_snapshot, created_at
		FROM memory
		WHERE plan_id = $1 AND next_message_id IS NULL
		LIMIT 1
	`

	var message models.Message
	if err := pgxscan.Get(ctx, q, &message, query, planID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}

	return &message, nil
}

// DeleteConversation deletes all messages for a given plan.
func (s *MemoryStore) DeleteConversation(ctx context.Context, planID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE plan_id = $1`, MemoryTableName)
	_, err := s.db.Exec(ctx, query, planID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

// DeleteMessage deletes a single message and repairs the linked list.
func (s *MemoryStore) DeleteMessage(ctx context.Context, messageID string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get the message to be deleted to find its neighbors
	var msg models.Message
	query := fmt.Sprintf(`SELECT previous_message_id, next_message_id FROM %s WHERE id = $1`, MemoryTableName)
	err = tx.QueryRow(ctx, query, messageID).Scan(&msg.PreviousMessageID, &msg.NextMessageID)
	if err != nil {
		return fmt.Errorf("failed to get message to delete: %w", err)
	}

	// Update previous message to point to next message
	if msg.PreviousMessageID != nil {
		updatePrev := fmt.Sprintf(`UPDATE %s SET next_message_id = $1 WHERE id = $2`, MemoryTableName)
		_, err = tx.Exec(ctx, updatePrev, msg.NextMessageID, *msg.PreviousMessageID)
		if err != nil {
			return fmt.Errorf("failed to update previous message: %w", err)
		}
	}

	// Update next message to point to previous message
	if msg.NextMessageID != nil {
		updateNext := fmt.Sprintf(`UPDATE %s SET previous_message_id = $1 WHERE id = $2`, MemoryTableName)
		_, err = tx.Exec(ctx, updateNext, msg.PreviousMessageID, *msg.NextMessageID)
		if err != nil {
			return fmt.Errorf("failed to update next message: %w", err)
		}
	}

	// Delete the message
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, MemoryTableName)
	_, err = tx.Exec(ctx, deleteQuery, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteMessagesAfter deletes the given message and all subsequent messages in the conversation.
func (s *MemoryStore) DeleteMessagesAfter(ctx context.Context, messageID string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get the message to find its previous message
	var msg models.Message
	query := fmt.Sprintf(`SELECT previous_message_id FROM %s WHERE id = $1`, MemoryTableName)
	err = tx.QueryRow(ctx, query, messageID).Scan(&msg.PreviousMessageID)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	// Recursive query to find all subsequent messages (including the target message)
	// We start with the target messageID
	deleteQuery := fmt.Sprintf(`
		WITH RECURSIVE chain AS (
			SELECT id, next_message_id
			FROM %s
			WHERE id = $1

			UNION ALL

			SELECT m.id, m.next_message_id
			FROM %s m
			INNER JOIN chain c ON m.previous_message_id = c.id
		)
		DELETE FROM %s WHERE id IN (SELECT id FROM chain)
	`, MemoryTableName, MemoryTableName, MemoryTableName)

	_, err = tx.Exec(ctx, deleteQuery, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete messages chain: %w", err)
	}

	// Update the previous message's next_message_id to NULL (since it's now the last message)
	if msg.PreviousMessageID != nil {
		updatePrev := fmt.Sprintf(`UPDATE %s SET next_message_id = NULL WHERE id = $1`, MemoryTableName)
		_, err = tx.Exec(ctx, updatePrev, *msg.PreviousMessageID)
		if err != nil {
			return fmt.Errorf("failed to update previous message: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateMessage updates the content and/or plan snapshot of a message.
func (s *MemoryStore) UpdateMessage(ctx context.Context, messageID, content string, planSnapshot *models.Plan) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET content = $1, plan_snapshot = $2
		WHERE id = $3
	`, MemoryTableName)
	_, err := s.db.Exec(ctx, query, content, planSnapshot, messageID)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}
