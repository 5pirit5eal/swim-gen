package models

import (
	"time"
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
