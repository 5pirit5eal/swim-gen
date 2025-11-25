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
	ID                string    `json:"id" db:"id"`
	PlanID            string    `json:"plan_id" db:"plan_id"`
	UserID            string    `json:"user_id" db:"user_id"`
	Role              Role      `json:"role" db:"role"`
	Content           string    `json:"content" db:"content"`
	PreviousMessageID *string   `json:"previous_message_id" db:"previous_message_id"`
	NextMessageID     *string   `json:"next_message_id" db:"next_message_id"`
	PlanSnapshot      *Plan     `json:"plan_snapshot,omitempty" db:"plan_snapshot"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
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
