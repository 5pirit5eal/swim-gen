package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/5pirit5eal/swim-rag/internal/rag"
	"github.com/5pirit5eal/swim-rag/internal/scraper"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/schema"
)

type RAGServer struct {
	ctx    context.Context
	db     *rag.RAGDB
	client llms.Model
	cfg    models.Config
}

// NewRAGServer initializes a new RAG server with the given configuration.
// It loads the database password from Google Secret Manager and initializes
// the database connection and LLM client.
// It returns a pointer to the RAGServer and an error if any occurred during
// initialization.
func NewRAGServer(ctx context.Context, cfg models.Config) (*RAGServer, error) {
	log.Println("Initializing RAG server with config:", cfg)
	// Initialize the LLM client
	client, err := googleai.New(
		ctx, googleai.WithCloudProject(cfg.ProjectID),
		googleai.WithCloudLocation(cfg.Region),
		googleai.WithDefaultModel(cfg.Model),
		googleai.WithDefaultEmbeddingModel(cfg.Embedding.Model),
		googleai.WithHarmThreshold(googleai.HarmBlockLowAndAbove),
		googleai.WithAPIKey(cfg.APIKey),
	)
	if err != nil {
		return nil, err
	}

	db, err := rag.NewStore(client, ctx, cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Creating database connection successfully")

	return &RAGServer{
		ctx:    ctx,
		client: client,
		cfg:    cfg,
		db:     db,
	}, nil
}

func (rs *RAGServer) AddDocuments(w http.ResponseWriter, req *http.Request) {
	log.Println("Adding documents to the database...")
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
	ids, err := rs.db.Store.AddDocuments(rs.ctx, documents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"success","ids":` + fmt.Sprintf("%v", ids) + `}`))
}

func (rs *RAGServer) Query(w http.ResponseWriter, req *http.Request) {
	log.Println("Querying the database...")
	// Parse HTTP request from JSON.

	qr := &models.QueryRequest{}
	err := models.GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer, err := rs.db.Query(rs.ctx, rs.client, qr.Content, qr.Filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Answer generated successfully")

	models.WriteResponseJSON(w, http.StatusOK, answer)
}

func (rs *RAGServer) Close() {
	if err := rs.db.Store.Close(); err != nil {
		log.Printf("error closing store: %v", err)
	}
}

func (rs *RAGServer) Scrape(w http.ResponseWriter, rq *http.Request) {
	log.Println("Getting scraping request...")
	// Parse the URL from the request
	url := rq.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	log.Println("...for URL:", url)

	// Scrape the URL
	err := scraper.ScrapeURL(rs.db, rs.ctx, rs.client, rs.cfg, url)
	if err != nil {
		http.Error(w, fmt.Errorf("failed to scrape URL %s: %w", url, err).Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
