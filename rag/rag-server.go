package rag

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

type RAGServer struct {
	ctx         context.Context
	conn        *pgxpool.Pool
	store       *pgvector.Store
	modelClient llms.Model
	config      Config
}

// GetSecret retrieves the database password from Google Secret Manager.
// It takes a context and the secret location as parameters.
// It returns the password as a string and an error if any occurred during
// retrieval.
// The secret location should be in the format:
// "projects/{project_id}/secrets/{secret_name}/versions/latest".
func GetSecret(ctx context.Context, location string) (string, error) {
	log.Printf("Getting %s from secret manager", location)
	// Create a new Secret Manager client
	// and access the secret version.
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer c.Close()
	secret, err := c.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: location,
	})
	if err != nil {
		return "", err
	}
	log.Println("Got DB password from secret manager successfully")
	// The secret payload is a byte array, so convert it to a string.
	return string(secret.Payload.Data), nil
}

// NewRAGServer initializes a new RAG server with the given configuration.
// It loads the database password from Google Secret Manager and initializes
// the database connection and LLM client.
// It returns a pointer to the RAGServer and an error if any occurred during
// initialization.
func NewRAGServer(ctx context.Context, config Config) (*RAGServer, error) {
	log.Println("Initializing RAG server with config:", config)
	// Initialize the LLM client
	vertexClient, err := vertex.New(
		ctx, googleai.WithCloudProject(config.ProjectID),
		googleai.WithCloudLocation(config.Region),
		googleai.WithDefaultModel(config.Model),
		googleai.WithDefaultEmbeddingModel(config.Embedding.Model),
		googleai.WithHarmThreshold(googleai.HarmBlockLowAndAbove),
	)
	if err != nil {
		return nil, err
	}

	// Load the database password from Google Secret Manager
	if config.DB.Pass == "" {
		pass, err := GetSecret(ctx, config.DB.PassLocation)
		if err != nil {
			return nil, err
		}
		config.DB.Pass = pass
	}
	log.Println("Got DB password successfully")

	// Create an embedder
	embedder, err := embeddings.NewEmbedder(vertexClient)
	if err != nil {
		return nil, err
	}

	log.Println("Creating database connection...")
	// Initialize the database connection
	// TODO: connect via cloud sql proxy (Unix socket or TCP) or directly via URL
	config.DB.Instance = fmt.Sprintf("host=/tmp/cloudsql/%s database=%s user=%s password=%s",
		config.DB.Instance, config.DB.Name, config.DB.User, config.DB.Pass)
	conn, err := pgxpool.New(ctx, config.DB.Instance)
	if err != nil {
		return nil, err
	}
	log.Println("Database connection created successfully")
	if _, err = conn.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		log.Println("Failed to create vector extension:", err.Error())
		return nil, err
	}
	log.Println("Created vector extension successfully")
	// Create a new store
	store, err := pgvector.New(
		ctx, pgvector.WithConn(conn),
		pgvector.WithEmbeddingTableName(config.Embedding.Name),
		pgvector.WithCollectionTableName("documents"),
		pgvector.WithEmbedder(embedder),
		pgvector.WithVectorDimensions(config.Embedding.Size),
	)
	if err != nil {
		return nil, err
	}
	log.Println("Created store successfully")
	// Create the URL table if it doesn't exist
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if err := createURLTableIfNotExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	log.Println("Creating database connection successfully")

	return &RAGServer{
		ctx:         ctx,
		store:       &store,
		modelClient: vertexClient,
		config:      config,
		conn:        conn,
	}, nil
}

func createURLTableIfNotExists(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946497); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create urls table: %w", err)
	}
	return nil
}

type document struct {
	Text     string         `json:"text"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
type addRequest struct {
	Documents []document `json:"documents"`
}

func (rs *RAGServer) AddDocuments(w http.ResponseWriter, req *http.Request) {
	log.Println("Adding documents to the database...")
	// Parse HTTP request from JSON.

	ar := &addRequest{}

	err := GetRequestJSON(req, ar)
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
	ids, err := rs.store.AddDocuments(rs.ctx, documents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","ids":` + fmt.Sprintf("%v", ids) + `}`))
}

func (rs *RAGServer) Query(w http.ResponseWriter, req *http.Request) {
	log.Println("Querying the database...")
	// Parse HTTP request from JSON.
	type queryRequest struct {
		Content string            `json:"content"`
		Filter  map[string]string `json:"filter,omitempty"`
	}
	qr := &queryRequest{}
	err := GetRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the most similar documents.
	docs, err := rs.store.SimilaritySearch(rs.ctx, qr.Content, 10, vectorstores.WithFilters(qr.Filter))
	if err != nil {
		http.Error(w, fmt.Errorf("similarity search: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	var docsContents []string
	for _, doc := range docs {
		docsContents = append(docsContents, doc.PageContent)
	}

	log.Printf("Found %d documents", len(docsContents))

	// Create a RAG query for the LLM with the most relevant documents as context
	query := fmt.Sprintf(ragTemplateStr, qr.Content, strings.Join(docsContents, "\n"))
	answer, err := llms.GenerateFromSinglePrompt(rs.ctx, rs.modelClient, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Answer generated successfully")

	WriteResponseJSON(w, http.StatusOK, answer)
}

func (rs *RAGServer) Close() {
	if err := rs.store.Close(); err != nil {
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

	// Load urls in the database into the scraper
	alreadyVisited := make([]string, 0)

	rows, err := rs.conn.Query(rs.ctx, "select url from urls")
	if err != nil {
		http.Error(w, "Failed to query database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Queried database successfully")
	defer rows.Close()
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			http.Error(w, "Failed to scan database row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		alreadyVisited = append(alreadyVisited, url)
	}

	// Scrape the URL
	plans, err := Scrape(alreadyVisited, url)
	if err != nil {
		http.Error(w, "Failed to scrape URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dc := make(chan schema.Document)
	ec := make(chan error)
	documents := make([]schema.Document, 0)
	errors := make([]error, 0)

	for kvp := range plans.Range() {
		go improvePlan(rs.ctx, rs.modelClient, kvp.Plan, dc, ec)
	}

	for range plans.Range() {
		select {
		case doc := <-dc:
			documents = append(documents, doc)
		case err := <-ec:
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		http.Error(w, fmt.Sprintf("Failed to improve plans: %v", errors), http.StatusInternalServerError)
		return
	}
	log.Println("Adding documents to the database...")

	// Store documents and their embeddings in the database
	ids, err := rs.store.AddDocuments(rs.ctx, documents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Added documents to the database successfully")

	batch := &pgx.Batch{}
	for kvp := range plans.Range() {
		log.Println("Inserting URL into database:", kvp.URL)
		batch.Queue("INSERT INTO urls (url) VALUES ($1) ON CONFLICT (url) DO NOTHING", kvp.URL)
	}
	br := rs.conn.SendBatch(rs.ctx, batch)

	if err := br.Close(); err != nil {
		http.Error(w, "Failed to insert urls into database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Inserted URLs into database successfully")

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	urls := make([]string, len(documents))
	for i, doc := range documents {
		urls[i] = doc.Metadata["url"].(string)
	}
	w.Write([]byte(`{"status":"success","ids":` + fmt.Sprintf("%v", ids) + `,"urls":` + fmt.Sprintf("%v", urls) + `}`))
}
