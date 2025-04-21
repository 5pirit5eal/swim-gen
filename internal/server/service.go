package server

import (
	"context"
	"net/http"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/5pirit5eal/swim-rag/internal/rag"
	"github.com/go-chi/httplog/v2"
	"github.com/tmc/langchaingo/schema"
)

type RAGService struct {
	// Background context for the server
	ctx context.Context
	// Database client used for storing and querying documents
	db *rag.RAGDB
	// Configuration for the RAG server
	cfg models.Config
}

// Initializes a new RAG service with the given configuration.
// It loads the database password from Google Secret Manager and initializes
// the database connection and LLM client.
// It returns a pointer to the RAGService and an error if any occurred during
// initialization.
func NewRAGService(ctx context.Context, cfg models.Config) (*RAGService, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("Initializing RAG server with config", "cfg", httplog.StructValue(cfg))
	db, err := rag.NewGoogleAIStore(ctx, cfg)
	if err != nil {
		return nil, err
	}

	logger.Info("Creating database connection successfully")

	return &RAGService{
		ctx: ctx,
		cfg: cfg,
		db:  db,
	}, nil
}

func (rs *RAGService) AddDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Adding documents to the database...")
	// Parse HTTP request from JSON.

	ar := &models.AddRequest{}

	err := models.GetRequestJSON(req, ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the documents to the format expected by the store
	var documents []schema.Document
	for _, doc := range ar.Documents {
		documents = append(documents, schema.Document{PageContent: doc.Text, Metadata: doc.Metadata})
	}

	// Store documents and their embeddings in the database
	ids, err := rs.db.Store.AddDocuments(req.Context(), documents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with a success message
	models.WriteResponseJSON(w, http.StatusOK, models.AddResponse{Status: "OK", IDs: ids})
}

func (rs *RAGService) QueryHandler(w http.ResponseWriter, req *http.Request) {
	logger := httplog.LogEntry(req.Context())
	logger.Info("Querying the database...")
	// Parse HTTP request from JSON.

	qr := &models.QueryRequest{}
	err := models.GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer, err := rs.db.Query(req.Context(), qr.Content, qr.Filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info("Answer generated successfully")
	models.WriteResponseJSON(w, http.StatusOK, answer)
}
