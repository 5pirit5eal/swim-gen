package models

import "time"

// UploadPlanRequest represents the request body for donating a training plan
// @Description Request payload for donating a swim training plan to the system
type UploadPlanRequest struct {
	Title        string   `json:"title,omitempty" example:"Advanced Freestyle Training"`
	Description  string   `json:"description,omitempty" example:"A comprehensive training plan for improving freestyle technique"`
	Table        Table    `json:"table" binding:"required"`
	Language     Language `json:"language,omitempty" example:"en"` // Language specifies the language of the training plan
	AllowSharing bool     `json:"allow_sharing"`                   // AllowSharing indicates if the plan can be shared with others
}

// QueryRequest represents the request body for querying the RAG system
// @Description Request payload for querying swim training plans from the RAG system
type QueryRequest struct {
	Content     string         `json:"content" example:"I need a training plan for improving my freestyle technique" binding:"required"` // Content describes what kind of training plan is needed
	Filter      map[string]any `json:"filter,omitempty"`                                                                                 // Filter allows filtering plans by metadata like difficulty or stroke type
	Method      string         `json:"method" example:"generate" validate:"oneof=choose generate" binding:"required"`                    // Method can be either 'choose' (select existing plan) or 'generate' (create new plan)
	Language    Language       `json:"language,omitempty" example:"en"`                                                                  // Language specifies the language for the response
	PoolLength  any            `json:"pool_length,omitempty" validate:"oneof=25 50 Freiwasser"`                                          // PoolLength specifies the pool length for the training plan
	Preferences *bool          `json:"preferences,omitempty"`                                                                            // Preferences indicates if the user profile should be used for generation
}

// RAGResponse represents the response after a query to the RAG system
// @Description Response containing a generated or selected swim training plan
type RAGResponse struct {
	PlanID      string `json:"plan_id,omitempty" example:"plan_123"` // PlanID is the identifier of the training plan
	Title       string `json:"title" example:"Advanced Freestyle Training"`
	Description string `json:"description" example:"A comprehensive training plan for improving freestyle technique"`
	Table       Table  `json:"table"`
}

func (r *RAGResponse) Plan() *Plan {
	if r == nil {
		return nil
	}
	return &Plan{
		PlanID:      r.PlanID,
		Title:       r.Title,
		Description: r.Description,
		Table:       r.Table,
	}
}

func (r *RAGResponse) Map() map[string]any {
	return map[string]any{
		"plan_id":     r.PlanID,
		"title":       r.Title,
		"description": r.Description,
	}
}

// PlanToPDFRequest represents the request for PDF export
// @Description Request payload for exporting a training plan to PDF format
type PlanToPDFRequest struct {
	PlanID          string   `json:"plan_id,omitempty" example:"plan_123"` // PlanID identifies the training plan to be exported
	Title           string   `json:"title" example:"Advanced Freestyle Training" binding:"required"`
	Description     string   `json:"description" example:"A comprehensive training plan for improving freestyle technique" binding:"required"`
	Table           Table    `json:"table" binding:"required"`
	Horizontal      bool     `json:"horizontal" example:"false"`                                 // Horizontal indicates if the PDF should be in landscape orientation
	LargeFont       bool     `json:"large_font" example:"true"`                                  // LargeFont indicates if the PDF should use a larger font size
	Language        Language `json:"language,omitempty" example:"en"`                            // Language specifies the language for the PDF content
	FrontendBaseURL string   `json:"frontend_base_url,omitempty" example:"https://swim-gen.app"` // FrontendBaseURL is the base URL for drill links in the PDF
}

// PlanToPDFResponse represents the response from PDF export
// @Description Response containing the URI to the generated PDF file
type PlanToPDFResponse struct {
	URI string `json:"uri" example:"https://storage.googleapis.com/bucket/plans/plan_123.pdf"`
}

// GeneratePromptRequest represents the request for prompt generation
// @Description Request payload for generating a prompt for swim training plan creation
type GeneratePromptRequest struct {
	Language Language `json:"language" example:"en" binding:"required"`
}

// GeneratedPromptResponse represents the response containing the generated prompt
// @Description Response containing the generated prompt for swim training plan creation
type GeneratedPromptResponse struct {
	Prompt string `json:"prompt" example:"Generate a swim training plan for improving freestyle technique"`
}

// HealthStatus represents the health status of the service and its components
// @Description Health status of the service and its components
type HealthStatus struct {
	Status        string            `json:"status"`
	Timestamp     string            `json:"timestamp"`
	Components    map[string]string `json:"components"`
	SchemaVersion int               `json:"schema_version,omitempty"`
}

// UpsertPlanRequest represents the request payload for upserting a training plan
// @Description Request payload for upserting a swim training plan to the system
type UpsertPlanRequest struct {
	PlanID      string `json:"plan_id,omitempty" example:"plan_123"` // PlanID identifies the training plan to be upserted
	Title       string `json:"title" example:"Advanced Freestyle Training" binding:"required"`
	Description string `json:"description" example:"A comprehensive training plan for improving freestyle technique" binding:"required"`
	Table       Table  `json:"table" binding:"required"`
}

// UpsertPlanResponse represents the response after upserting a training plan
// @Description Response containing the upserted swim training plan
type UpsertPlanResponse struct {
	PlanID string `json:"plan_id" example:"plan_123"`
}

// SharePlanRequest represents the request payload for sharing a training plan
// @Description Request payload for sharing a swim training plan
type SharePlanRequest struct {
	PlanID string        `json:"plan_id" example:"plan_123" binding:"required"` // PlanID identifies the training plan to be shared
	Method SharingMethod `json:"method" example:"link" binding:"required"`      // Method specifies the sharing method, e.g., 'link' or 'email'
}

// SharePlanResponse represents the response after sharing a training plan
// @Description Response containing the sharing details of the swim training plan
type SharePlanResponse struct {
	URLHash string `json:"url_hash" example:"abc123"` // URLHash is the hash to access the shared training plan
}

// ChatRequest represents the request payload for chat-based plan refinement
// @Description Request payload for conversational training plan creation and refinement
type ChatRequest struct {
	PlanID     string   `json:"plan_id,omitempty" example:"plan_123"`                          // PlanID identifies the conversation/plan (optional for new conversations)
	Message    string   `json:"message" example:"Make it more challenging" binding:"required"` // Message is the user's input to the chat
	Language   Language `json:"language,omitempty" example:"en"`                               // Language specifies the language for the response
	PoolLength any      `json:"pool_length,omitempty" validate:"oneof=25 50 Freiwasser"`       // PoolLength specifies the pool length for the training plan
}

// ChatResponsePayload represents the response from a chat interaction
// @Description Response containing the updated plan and conversational response
type ChatResponsePayload struct {
	PlanID      string `json:"plan_id" example:"plan_123"`                                                                      // PlanID identifies the conversation/plan
	Title       string `json:"title,omitempty" example:"Advanced Freestyle Training"`                                           // Title of the training plan
	Description string `json:"description,omitempty" example:"A comprehensive training plan for improving freestyle technique"` // Description of the training plan
	Table       Table  `json:"table,omitempty"`                                                                                 // Table containing the training plan details
	Response    string `json:"response" example:"I've made the plan more challenging by adding butterfly sets"`                 // Response is the conversational AI response explaining changes
}

// PlanSnapshot represents a snapshot of a training plan
// @Description Snapshot of a training plan
type MessagePayload struct {
	ID                string       `json:"id" db:"id"`
	PlanID            string       `json:"plan_id" db:"plan_id"`
	UserID            string       `json:"user_id" db:"user_id"`
	Role              Role         `json:"role" db:"role"`
	Content           string       `json:"content" db:"content"`
	PreviousMessageID *string      `json:"previous_message_id" db:"previous_message_id"`
	NextMessageID     *string      `json:"next_message_id" db:"next_message_id"`
	PlanSnapshot      *RAGResponse `json:"plan_snapshot" db:"plan_snapshot"`
	CreatedAt         time.Time    `json:"created_at" db:"created_at"` // Table containing the training plan details
}

// GetConversationResponse represents the response from a conversation history request
// @Description Response containing the conversation history
type GetConversationResponse struct {
	Conversation []MessagePayload `json:"conversation"` // Conversation history
}

// DeleteMessageRequest represents the request payload for deleting a single message
// @Description Request payload for deleting a single message from conversation history
type DeleteMessageRequest struct {
	MessageID string `json:"message_id" example:"msg_123" binding:"required"` // MessageID identifies the message to delete
}

// DeleteMessagesAfterRequest represents the request payload for deleting a message and all subsequent messages
// @Description Request payload for deleting a message and all subsequent messages in the conversation
type DeleteMessagesAfterRequest struct {
	MessageID string `json:"message_id" example:"msg_123" binding:"required"` // MessageID identifies the message from which to delete (inclusive)
}

// DeleteConversationRequest represents the request payload for deleting an entire conversation
// @Description Request payload for deleting an entire conversation and all its messages
type DeleteConversationRequest struct {
	PlanID string `json:"plan_id" example:"plan_123" binding:"required"` // PlanID identifies the conversation to delete
}

// AddMessageRequest represents the request body for adding a message to the conversation history
// @Description Request payload for adding a message to the conversation history
type AddMessageRequest struct {
	PlanID            string       `json:"plan_id" binding:"required"`
	Role              Role         `json:"role" binding:"required"`
	Content           string       `json:"content" binding:"required"`
	PreviousMessageID string       `json:"previous_message_id,omitempty"`
	PlanSnapshot      *RAGResponse `json:"plan_snapshot,omitempty"`
}

// AddPlanToHistoryRequest represents the request payload for adding a plan to history
// @Description Request payload for adding a plan to the authenticated user's history
type AddPlanToHistoryRequest struct {
	PlanID      string `json:"plan_id" example:"plan_123" binding:"required"`                                                            // PlanID identifies the plan to add to history
	Title       string `json:"title" example:"Advanced Freestyle Training" binding:"required"`                                           // Title of the plan
	Description string `json:"description" example:"A comprehensive training plan for improving freestyle technique" binding:"required"` // Description of the plan
	Table       Table  `json:"table" binding:"required"`                                                                                 // Table containing the plan details
}

func (a *AddPlanToHistoryRequest) Plan() *Plan {
	return &Plan{
		PlanID:      a.PlanID,
		Title:       a.Title,
		Description: a.Description,
		Table:       a.Table,
	}
}

// AddPlanToHistoryResponse represents the response after adding a plan to history
// @Description Response containing the new plan ID and a success message
type AddPlanToHistoryResponse struct {
	Message string `json:"message" example:"Plan added to history successfully"`
	PlanID  string `json:"plan_id" example:"plan_123"`
}

// FeedbackRequest represents the request payload for submitting feedback
// @Description Request payload for submitting feedback on a training plan
type FeedbackRequest struct {
	PlanID           string `json:"plan_id" example:"plan_123" binding:"required"`
	Rating           int    `json:"rating" example:"5" binding:"required" validate:"min=1,max=5"`
	WasSwam          bool   `json:"was_swam" example:"true"`
	DifficultyRating int    `json:"difficulty_rating" example:"7" validate:"min=1,max=10"`
	Comment          string `json:"comment,omitempty" example:"Great plan!"`
}
