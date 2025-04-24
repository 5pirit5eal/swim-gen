package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/5pirit5eal/swim-rag/internal/config"
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
	cfg config.Config
}

// Initializes a new RAG service with the given configuration.
// It loads the database password from Google Secret Manager and initializes
// the database connection and LLM client.
// It returns a pointer to the RAGService and an error if any occurred during
// initialization.
func NewRAGService(ctx context.Context, cfg config.Config) (*RAGService, error) {
	slog.Info("Initializing RAG server with config", "cfg", slog.AnyValue(cfg))
	db, err := rag.NewGoogleAIStore(ctx, cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("Creating database connection successfully")

	return &RAGService{
		ctx: ctx,
		cfg: cfg,
		db:  db,
	}, nil
}

// Closes the database connection and LLM client.
// It is important to call this method when the service is no longer needed
// to release resources and avoid memory leaks.
func (rs *RAGService) Close() {
	slog.Info("Closing RAG server...")
	if err := rs.db.Store.Close(); err != nil {
		slog.Error("Error closing database connection", "err", err.Error())
	}
	slog.Info("RAG server closed successfully")
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

	// Recalculate the sums of the rows to be sure they are correct
	answer.Table.UpdateSum()

	logger.Info("Answer generated successfully")
	models.WriteResponseJSON(w, http.StatusOK, answer)
}
