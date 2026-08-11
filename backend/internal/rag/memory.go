package rag

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

var (
	ErrMemoryNotFound   = errors.New("memory resource not found")
	ErrMemoryValidation = errors.New("invalid memory request")
)

type MemoryStore struct {
	db *pgxpool.Pool
}

// Ensure MemoryStore implements Memory
var _ models.Memory = (*MemoryStore)(nil)

func NewMemoryStore(db *pgxpool.Pool) *MemoryStore {
	return &MemoryStore{db: db}
}

func (s *MemoryStore) userOwnsPlan(ctx context.Context, q pgxscan.Querier, planID, userID string) (bool, error) {
	var ownsPlan bool
	err := pgxscan.Get(ctx, q, &ownsPlan, fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1 FROM %s WHERE plan_id = $1 AND user_id = $2
			UNION ALL
			SELECT 1 FROM %s WHERE plan_id = $1 AND user_id = $2
		)
	`, HistoryTableName, DonatedPlanTable), planID, userID)
	return ownsPlan, err
}

// AddMessage inserts a new message into the memory table and updates the linked list.
func (s *MemoryStore) AddMessage(ctx context.Context, planID, userID string, role models.Role, content string, previousMessageID *string, planSnapshot *models.Plan) (*models.Message, error) {
	if planSnapshot != nil && planSnapshot.PlanID != planID {
		return nil, fmt.Errorf("%w: plan snapshot does not match plan", ErrMemoryValidation)
	}
	if role != models.RoleUser && role != models.RoleAI {
		return nil, fmt.Errorf("%w: unsupported message role", ErrMemoryValidation)
	}
	if planID == "" || userID == "" || content == "" {
		return nil, fmt.Errorf("%w: plan, user, and content are required", ErrMemoryValidation)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	ownsPlan, err := s.userOwnsPlan(ctx, tx, planID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify plan ownership: %w", err)
	}
	if !ownsPlan {
		return nil, ErrMemoryNotFound
	}

	// If PreviousMessageID is not provided, try to find the last message in the conversation
	if previousMessageID == nil {
		lastMsg, err := s.GetLastMessage(ctx, tx, planID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get last message: %w", err)
		}
		if lastMsg != nil {
			previousMessageID = &lastMsg.ID
		}
	} else {
		var exists bool
		if err := tx.QueryRow(ctx, fmt.Sprintf(`
			SELECT EXISTS (
				SELECT 1 FROM %s
				WHERE id = $1 AND plan_id = $2 AND user_id = $3
			)
		`, MemoryTableName), *previousMessageID, planID, userID).Scan(&exists); err != nil {
			return nil, fmt.Errorf("failed to verify previous message: %w", err)
		}
		if !exists {
			return nil, ErrMemoryNotFound
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
			WHERE id = $2 AND plan_id = $3 AND user_id = $4
		`, MemoryTableName)
		result, err := tx.Exec(ctx, updateQuery, newMessageID, *previousMessageID, planID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to update previous message: %w", err)
		}
		if result.RowsAffected() != 1 {
			return nil, ErrMemoryNotFound
		}
	}

	// If a valid plan snapshot is provided, update the plan table
	if planSnapshot != nil {
		upsertPlanQuery := fmt.Sprintf(`
			INSERT INTO %s (plan_id, title, description, plan_table)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (plan_id) DO UPDATE
			SET title = EXCLUDED.title,
				description = EXCLUDED.description,
				plan_table = EXCLUDED.plan_table,
				updated_at = now()
		`, PlanTableName)
		_, err = tx.Exec(ctx, upsertPlanQuery, planID, planSnapshot.Title, planSnapshot.Description, planSnapshot.Table)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert plan: %w", err)
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
func (s *MemoryStore) GetConversation(ctx context.Context, planID, userID string) ([]models.Message, error) {
	ownsPlan, err := s.userOwnsPlan(ctx, s.db, planID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify plan ownership: %w", err)
	}
	if !ownsPlan {
		return nil, ErrMemoryNotFound
	}

	query := `
		WITH RECURSIVE conversation AS (
			SELECT id, plan_id, user_id, role, content, previous_message_id, next_message_id, plan_snapshot, created_at, 0 AS depth, ARRAY[id] AS path
			FROM memory
			WHERE plan_id = $1 AND user_id = $2 AND previous_message_id IS NULL

			UNION ALL

			SELECT m.id, m.plan_id, m.user_id, m.role, m.content, m.previous_message_id, m.next_message_id, m.plan_snapshot, m.created_at, c.depth + 1, c.path || m.id
			FROM memory m
			INNER JOIN conversation c ON m.previous_message_id = c.id
				AND m.plan_id = c.plan_id
				AND m.user_id = c.user_id
			WHERE m.plan_id = $1 AND m.user_id = $2 AND NOT m.id = ANY(c.path)
		)
		SELECT id, plan_id, user_id, role, content, previous_message_id, next_message_id, plan_snapshot, created_at
		FROM conversation
		ORDER BY depth;
	`

	messages := make([]models.Message, 0)
	if err := pgxscan.Select(ctx, s.db, &messages, query, planID, userID); err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	return messages, nil
}

// GetLastMessage retrieves the last message in the conversation (where next_message_id is NULL).
// It accepts a querier (tx or pool) to support transactions.
func (s *MemoryStore) GetLastMessage(ctx context.Context, q pgxscan.Querier, planID, userID string) (*models.Message, error) {
	query := `
		SELECT id, plan_id, user_id, role, content, previous_message_id, next_message_id, plan_snapshot, created_at
		FROM memory
		WHERE plan_id = $1 AND user_id = $2 AND next_message_id IS NULL
		LIMIT 1
	`

	var message models.Message
	if err := pgxscan.Get(ctx, q, &message, query, planID, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last message: %w", err)
	}

	return &message, nil
}

// DeleteConversation deletes all messages for a plan owned by the user.
func (s *MemoryStore) DeleteConversation(ctx context.Context, planID, userID string) error {
	ownsPlan, err := s.userOwnsPlan(ctx, s.db, planID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify plan ownership: %w", err)
	}
	if !ownsPlan {
		return ErrMemoryNotFound
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE plan_id = $1 AND user_id = $2`, MemoryTableName)
	_, err = s.db.Exec(ctx, query, planID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

// DeleteMessage deletes a single message and repairs the linked list.
func (s *MemoryStore) DeleteMessage(ctx context.Context, messageID, userID string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get the message to be deleted to find its neighbors
	var msg models.Message
	query := fmt.Sprintf(`
		SELECT plan_id, user_id, previous_message_id, next_message_id
		FROM %s WHERE id = $1 AND user_id = $2
	`, MemoryTableName)
	err = tx.QueryRow(ctx, query, messageID, userID).Scan(&msg.PlanID, &msg.UserID, &msg.PreviousMessageID, &msg.NextMessageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMemoryNotFound
		}
		return fmt.Errorf("failed to get message to delete: %w", err)
	}

	// Update previous message to point to next message
	if msg.PreviousMessageID != nil {
		updatePrev := fmt.Sprintf(`UPDATE %s SET next_message_id = $1 WHERE id = $2 AND plan_id = $3 AND user_id = $4`, MemoryTableName)
		result, err := tx.Exec(ctx, updatePrev, msg.NextMessageID, *msg.PreviousMessageID, msg.PlanID, userID)
		if err != nil {
			return fmt.Errorf("failed to update previous message: %w", err)
		}
		if result.RowsAffected() != 1 {
			return fmt.Errorf("failed to update previous message: %w", ErrMemoryNotFound)
		}
	}

	// Update next message to point to previous message
	if msg.NextMessageID != nil {
		updateNext := fmt.Sprintf(`UPDATE %s SET previous_message_id = $1 WHERE id = $2 AND plan_id = $3 AND user_id = $4`, MemoryTableName)
		result, err := tx.Exec(ctx, updateNext, msg.PreviousMessageID, *msg.NextMessageID, msg.PlanID, userID)
		if err != nil {
			return fmt.Errorf("failed to update next message: %w", err)
		}
		if result.RowsAffected() != 1 {
			return fmt.Errorf("failed to update next message: %w", ErrMemoryNotFound)
		}
	}

	// Delete the message
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND plan_id = $2 AND user_id = $3`, MemoryTableName)
	result, err := tx.Exec(ctx, deleteQuery, messageID, msg.PlanID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	if result.RowsAffected() != 1 {
		return ErrMemoryNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteMessagesAfter deletes the given message and all subsequent messages in the conversation.
func (s *MemoryStore) DeleteMessagesAfter(ctx context.Context, messageID, userID string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get the message to find its previous message
	var msg models.Message
	query := fmt.Sprintf(`
		SELECT plan_id, user_id, previous_message_id
		FROM %s WHERE id = $1 AND user_id = $2
	`, MemoryTableName)
	err = tx.QueryRow(ctx, query, messageID, userID).Scan(&msg.PlanID, &msg.UserID, &msg.PreviousMessageID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMemoryNotFound
		}
		return fmt.Errorf("failed to get message: %w", err)
	}
	if msg.PreviousMessageID != nil {
		var predecessorExists bool
		if err := tx.QueryRow(ctx, fmt.Sprintf(`
			SELECT EXISTS (
				SELECT 1 FROM %s
				WHERE id = $1 AND plan_id = $2 AND user_id = $3
			)
		`, MemoryTableName), *msg.PreviousMessageID, msg.PlanID, userID).Scan(&predecessorExists); err != nil {
			return fmt.Errorf("failed to verify previous message: %w", err)
		}
		if !predecessorExists {
			return ErrMemoryNotFound
		}
	}

	// Recursive query to find all subsequent messages (including the target message).
	// The path prevents malformed cycles from causing an unbounded query.
	deleteQuery := fmt.Sprintf(`
		WITH RECURSIVE chain AS (
			SELECT id, next_message_id, ARRAY[id] AS path
			FROM %s
			WHERE id = $1 AND plan_id = $2 AND user_id = $3

			UNION ALL

			SELECT m.id, m.next_message_id, c.path || m.id
			FROM %s m
			INNER JOIN chain c ON m.previous_message_id = c.id
			WHERE m.plan_id = $2 AND m.user_id = $3 AND NOT m.id = ANY(c.path)
		)

		DELETE FROM %s
		WHERE plan_id = $2 AND user_id = $3 AND id IN (SELECT id FROM chain)
	`, MemoryTableName, MemoryTableName, MemoryTableName)

	if msg.PreviousMessageID != nil {
		updatePrev := fmt.Sprintf(`UPDATE %s SET next_message_id = NULL WHERE id = $1 AND plan_id = $2 AND user_id = $3`, MemoryTableName)
		result, err := tx.Exec(ctx, updatePrev, *msg.PreviousMessageID, msg.PlanID, userID)
		if err != nil {
			return fmt.Errorf("failed to update previous message: %w", err)
		}
		if result.RowsAffected() != 1 {
			return fmt.Errorf("failed to update previous message: %w", ErrMemoryNotFound)
		}
	}

	result, err := tx.Exec(ctx, deleteQuery, messageID, msg.PlanID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete messages chain: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrMemoryNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
