package models

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

type Role string

const (
	RoleUser Role = "user"
	RoleAI   Role = "ai"
)

type Message struct {
	ID                string    `db:"id"`
	PlanID            string    `db:"plan_id"`
	UserID            string    `db:"user_id"`
	Role              Role      `db:"role"`
	Content           string    `db:"content"`
	PreviousMessageID *string   `db:"previous_message_id"`
	NextMessageID     *string   `db:"next_message_id"`
	PlanSnapshot      *Plan     `db:"plan_snapshot"`
	CreatedAt         time.Time `db:"created_at"`
}

type Memory interface {
	AddMessage(ctx context.Context, planID, userID string, role Role, content string, previousMessageID *string, planSnapshot *Plan) (*Message, error)
	GetConversation(ctx context.Context, planID string) ([]Message, error)
	GetLastMessage(ctx context.Context, q pgxscan.Querier, planID string) (*Message, error)
	DeleteConversation(ctx context.Context, planID string) error
	DeleteMessage(ctx context.Context, messageID string) error
	UpdateMessage(ctx context.Context, messageID, content string, planSnapshot *Plan) error
	DeleteMessagesAfter(ctx context.Context, messageID string) error
}
