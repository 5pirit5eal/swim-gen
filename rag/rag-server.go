package rag

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	pgx "github.com/jackc/pgx/v5"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

type Config struct {
	ProjectID string `env:"PROJECT_ID"`
	Region    string `env:"REGION"`
	Model     string `env:"MODEL"`
	Embedding struct {
		Name  string `env:"EMBEDDING_NAME"`
		Model string `env:"EMBEDDING_MODEL"`
		SIZE  int    `env:"EMBEDDING_SIZE"`
	}

	DB struct {
		Name         string `env:"DB_NAME"`
		IP           string `env:"DB_IP"`
		Port         string `env:"DB_PORT"`
		User         string `env:"DB_USER"`
		PassLocation string `env:"DB_PASS_LOCATION"`
		Pass         string
	}
}

type RAGServer struct {
	ctx         context.Context
	store       pgvector.Store
	modelClient llms.Model
	config      Config
}

// GetDBPass retrieves the database password from Google Secret Manager.
// It takes a context and the secret location as parameters.
// It returns the password as a string and an error if any occurred during
// retrieval.
// The secret location should be in the format:
// "projects/{project_id}/secrets/{secret_name}/versions/latest".
func GetDBPass(ctx context.Context, location string) (string, error) {
	log.Println("Getting DB password from secret manager")
	// Create a new Secret Manager client
	// and access the secret version.
	c, err := secretmanager.NewClient(ctx)
	defer c.Close()
	if err != nil {
		return "", err
	}
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
//
// Example usage:
//
//	ctx := context.Background()
//	config := Config{
//		ProjectID: "your-project-id",
//		Region:    "us-central1",
//		Model:     "your-model-name",
//		Embedding: struct {
//			Name:  "your-embedding-name",
//			Model: "your-embedding-model",
//			SIZE:  768,
//		},
//		DB: struct {
//			Name:         "your-db-name",
//			IP:           "your-db-ip",
//			Port:         "your-db-port",
//			User
//			PassLocation: "projects/your-project-id/secrets/your-secret-name/versions/latest",
//		},
//	}
//	ragServer, err := NewRAGServer(ctx, config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Use ragServer for further operations...
//	// ...
//	// Don't forget to close the server when done
//	defer ragServer.store.Close()
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
	pass, err := GetDBPass(ctx, config.DB.PassLocation)
	if err != nil {
		return nil, err
	}
	config.DB.Pass = pass
	log.Println("Got DB password successfully")

	// Create an embedder
	embedder, err := embeddings.NewEmbedder(vertexClient)
	if err != nil {
		return nil, err
	}

	log.Println("Creating database connection...")
	// Initialize the database connection
	conn, err := pgx.Connect(ctx, "postgres://"+config.DB.User+":"+pass+"@"+config.DB.IP+":"+config.DB.Port+"/"+config.DB.Name)
	store, err := pgvector.New(
		ctx, pgvector.WithConn(conn),
		pgvector.WithEmbeddingTableName(config.Embedding.Model+"-"+config.Embedding.Name),
		pgvector.WithCollectionTableName("documents"),
		pgvector.WithEmbedder(embedder),
		pgvector.WithVectorDimensions(config.Embedding.SIZE),
	)
	if err != nil {
		return nil, err
	}
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
		store:       store,
		modelClient: vertexClient,
		config:      config,
	}, nil
}

func createURLTableIfNotExists(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946497); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL UNIQUE
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

// TODO: Improve this augmentation prompt and make it more swimming specific
const ragTemplateStr = `
I will ask you a question and will provide some additional context information.
Assume this context information is factual and correct, as part of internal
documentation.
If the question relates to the context, answer it using the context.
If the question does not relate to the context, answer it as normal.

For example, let's say the context has nothing in it about tropical flowers;
then if I ask you about tropical flowers, just answer what you know about them
without referring to the context.

For example, if the context does mention minerology and I ask you about that,
provide information from the context along with general knowledge.

Question:
%s

Context:
%s
`

func (rs *RAGServer) ScrapeHandler(w http.ResponseWriter, rq *http.Request) {
	// Parse the URL from the request
	url := rq.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	// Load urls in the database into the scraper
	alreadyVisited := make([]string, 0)
	conn, err := pgx.Connect(rs.ctx, "postgres://"+rs.config.DB.User+":"+rs.config.DB.Pass+"@"+rs.config.DB.IP+":"+rs.config.DB.Port+"/"+rs.config.DB.Name)

	if err != nil {
		http.Error(w, "Failed to connect to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close(rs.ctx)
	rows, err := conn.Query(rs.ctx, "SELECT url FROM urls")
	if err != nil {
		http.Error(w, "Failed to query database: "+err.Error(), http.StatusInternalServerError)
		return
	}
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

	// Add the scraped data to a cloud task
	client, err := cloudtasks.NewClient(rs.ctx)
	if err != nil {
		http.Error(w, "Failed to create Cloud Tasks client: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Close()

	// TODO: Create request body by converting the plans into documents
	documents := make([]document, 0)
	for kvp := range plans.Range() {
		documents = append(documents, document{
			Text:     kvp.Plan.String(),
			Metadata: kvp.Plan.Map(),
		})
	}

	// TODO: Enhance scraped documents with gemini and create meaningful metadata
	// TODO: Write enhanced documents into db
	// TODO: Write newly scraped urls into db
	// TODO: Return Status OK

}
