package models

type Document struct {
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type AddRequest struct {
	Documents []Document `json:"documents"`
}

type AddResponse struct {
	Status string   `json:"status"`
	IDs    []string `json:"ids"`
}

type QueryRequest struct {
	Content string `json:"content"`
	// This has to be a map[string]any to support langchaingo filter syntax
	// pgvector currently only supports key=value pairs though
	Filter map[string]any `json:"filter,omitempty"`
}

type RAGResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Table       Table  `json:"table"`
}

type ChooseResponse struct {
	Idx         int    `json:"index"`
	Description string `json:"description"`
}

type PlanToPDFRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Table       Table  `json:"table"`
}

type PlanToPDFResponse struct {
	URI string `json:"uri"`
}
