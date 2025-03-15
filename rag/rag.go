package rag

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/golobby/dotenv"
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
	}
}

type RAGServer struct {
	ctx         context.Context
	store       pgvector.Store
	modelClient llms.Model
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

	// Create an embedder
	embedder, err := embeddings.NewEmbedder(vertexClient)
	if err != nil {
		return nil, err
	}

	log.Println("Creating database connection...")
	// Initialize the database connection
	store, err := pgvector.New(
		ctx, pgvector.WithConnectionURL("postgres://"+config.DB.User+":"+pass+"@"+config.DB.IP+":"+config.DB.Port+"/"+config.DB.Name),
		pgvector.WithEmbeddingTableName(config.Embedding.Model+"-"+config.Embedding.Name),
		pgvector.WithCollectionTableName("documents"),
		pgvector.WithEmbedder(embedder),
		pgvector.WithVectorDimensions(config.Embedding.SIZE),
	)
	if err != nil {
		return nil, err
	}

	log.Println("Creating database connection successfully")

	return &RAGServer{
		ctx:         ctx,
		store:       store,
		modelClient: vertexClient,
	}, nil
}

func (rs *RAGServer) AddDocuments(w http.ResponseWriter, req *http.Request) {
	log.Println("Adding documents to the database...")
	// Parse HTTP request from JSON.
	type document struct {
		Text     string         `json:"text"`
		Metadata map[string]any `json:"metadata,omitempty"`
	}
	type addRequest struct {
		Documents []document `json:"documents"`
	}
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
	docs, err := rs.store.SimilaritySearch(rs.ctx, qr.Content, 5, vectorstores.WithFilters(qr.Filter))
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

func (rs *RAGServer) Close() {
	if err := rs.store.Close(); err != nil {
		log.Printf("error closing store: %v", err)
	}
}

func Example(ctx context.Context, config Config) string {
	log.Println("Initializing example...")
	file, err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	if err := dotenv.NewDecoder(file).Decode(&config); err != nil {
		log.Fatal(err)
	}

	llm, err := vertex.New(ctx, googleai.WithCloudProject(config.ProjectID), googleai.WithCloudLocation(config.Region), googleai.WithDefaultModel(config.Model))
	if err != nil {
		log.Fatal(err)
	}

	prompt := "Tell me a danish joke in danish"
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Example generated successfully")

	return answer

}
