package rag

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/golobby/dotenv"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
	"github.com/tmc/langchaingo/schema"
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
	// Initialize the LLM client
	vertexClient, err := vertex.New(ctx, googleai.WithCloudProject(config.ProjectID), googleai.WithCloudLocation(config.Region), googleai.WithDefaultModel(config.Model), googleai.WithDefaultEmbeddingModel(config.Embedding.Model), googleai.WithHarmThreshold(googleai.HarmBlockLowAndAbove))
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

	// Initialize the database connection
	store, err := pgvector.New(
		ctx, pgvector.WithConnectionURL("postgres://"+config.DB.User+":"+pass+"@"+config.DB.IP+":"+config.DB.Port+"/"+config.DB.Name),
		pgvector.WithEmbeddingTableName(config.Embedding.Model+"-"+config.Embedding.Name),
		pgvector.WithCollectionTableName("documents"),
		pgvector.WithEmbedder(embedder),
	)
	if err != nil {
		return nil, err
	}

	return &RAGServer{
		ctx:         ctx,
		store:       store,
		modelClient: vertexClient,
	}, nil
}

func (rs *RAGServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type document struct {
		Text string `json:"text"`
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
		documents = append(documents, schema.Document{PageContent: doc.Text})
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

func Example(ctx context.Context, config Config) string {
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

	prompt := "What is the fastest swimming technique?"
	answer, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}
	return answer

}
