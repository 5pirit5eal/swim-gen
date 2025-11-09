package models

// DonatePlanRequest represents the request body for donating a training plan
// @Description Request payload for donating a swim training plan to the system
type DonatePlanRequest struct {
	Title       string   `json:"title,omitempty" example:"Advanced Freestyle Training"`
	Description string   `json:"description,omitempty" example:"A comprehensive training plan for improving freestyle technique"`
	Table       Table    `json:"table" binding:"required"`
	Language    Language `json:"language,omitempty" example:"en"` // Language specifies the language of the training plan
	// v3: add other table modalities
	// Image 	 string `json:"image,omitempty"`
	// URI 		 string `json:"uri,omitempty"`
}

// QueryRequest represents the request body for querying the RAG system
// @Description Request payload for querying swim training plans from the RAG system
type QueryRequest struct {
	Content    string         `json:"content" example:"I need a training plan for improving my freestyle technique" binding:"required"` // Content describes what kind of training plan is needed
	Filter     map[string]any `json:"filter,omitempty"`                                                                                 // Filter allows filtering plans by metadata like difficulty or stroke type
	Method     string         `json:"method" example:"generate" validate:"oneof=choose generate" binding:"required"`                    // Method can be either 'choose' (select existing plan) or 'generate' (create new plan)
	Language   Language       `json:"language,omitempty" example:"en"`                                                                  // Language specifies the language for the response
	PoolLength any            `json:"pool_length,omitempty" validate:"oneof=25 50 Freiwasser"`                                          // PoolLength specifies the pool length for the training plan
}

// RAGResponse represents the response after a query to the RAG system
// @Description Response containing a generated or selected swim training plan
type RAGResponse struct {
	PlanID      string `json:"plan_id,omitempty" example:"plan_123"` // PlanID is the identifier of the training plan
	Title       string `json:"title" example:"Advanced Freestyle Training"`
	Description string `json:"description" example:"A comprehensive training plan for improving freestyle technique"`
	Table       Table  `json:"table"`
}

// PlanToPDFRequest represents the request for PDF export
// @Description Request payload for exporting a training plan to PDF format
type PlanToPDFRequest struct {
	PlanID      string   `json:"plan_id,omitempty" example:"plan_123"` // PlanID identifies the training plan to be exported
	Title       string   `json:"title" example:"Advanced Freestyle Training" binding:"required"`
	Description string   `json:"description" example:"A comprehensive training plan for improving freestyle technique" binding:"required"`
	Table       Table    `json:"table" binding:"required"`
	Horizontal  bool     `json:"horizontal" example:"false"`      // Horizontal indicates if the PDF should be in landscape orientation
	LargeFont   bool     `json:"large_font" example:"true"`       // LargeFont indicates if the PDF should use a larger font size
	Language    Language `json:"language,omitempty" example:"en"` // Language specifies the language for the PDF content
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
