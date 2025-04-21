package rag

import (
	"context"
	"fmt"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/5pirit5eal/swim-rag/internal/config"
	"github.com/go-chi/httplog/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

const (
	CollectionTableName string = "documents"
)

type RAGDB struct {
	Conn   pgvector.PGXConn
	Store  *pgvector.Store
	Client llms.Model
}

func NewGoogleAIStore(ctx context.Context, cfg config.Config) (*RAGDB, error) {
	logger := httplog.LogEntry(ctx)
	// Initialize the LLM client
	client, err := googleai.New(
		ctx, googleai.WithCloudProject(cfg.ProjectID),
		googleai.WithCloudLocation(cfg.Region),
		googleai.WithDefaultModel(cfg.Model),
		googleai.WithDefaultEmbeddingModel(cfg.Embedding.Model),
		googleai.WithHarmThreshold(googleai.HarmBlockLowAndAbove),
		googleai.WithAPIKey(cfg.APIKey),
		googleai.WithDefaultMaxTokens(10000),
	)
	if err != nil {
		return nil, err
	}
	// Load the database password from Google Secret Manager
	if cfg.DB.Pass == "" {
		pass, err := GetSecret(ctx, cfg.DB.PassLocation)
		if err != nil {
			return nil, err
		}
		cfg.DB.Pass = pass
	}
	logger.Info("Got DB password successfully")

	// Create an embedder
	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}

	logger.Info("Creating database connection...")
	// Initialize the database connection
	conn, err := connect(ctx, cfg)
	if err != nil {
		return nil, err
	}
	logger.Info("Database connection created successfully")

	// Create a new store
	store, err := pgvector.New(
		ctx, pgvector.WithConn(conn),
		pgvector.WithEmbeddingTableName(cfg.Embedding.Name),
		pgvector.WithCollectionTableName(CollectionTableName),
		// Separate the collections and documents by embedding model name
		pgvector.WithCollectionName(cfg.Embedding.Model),
		pgvector.WithEmbedder(embedder),
		pgvector.WithVectorDimensions(cfg.Embedding.Size),
	)
	if err != nil {
		return nil, err
	}
	logger.Info("Created store successfully")
	// Create the URL table if it doesn't exist
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	if err := createURLTableIfNotExists(ctx, tx, cfg.Embedding.Name); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	logger.Info("Setup URL table successfully")
	return &RAGDB{Store: &store, Conn: conn, Client: client}, nil
}

func (rag *RAGDB) Close() error {
	if err := rag.Store.Close(); err != nil {
		return err
	}
	return nil
}

// GetSecret retrieves the database password from Google Secret Manager.
// It takes a context and the secret location as parameters.
// It returns the password as a string and an error if any occurred during
// retrieval.
// The secret location should be in the format:
// "projects/{project_id}/secrets/{secret_name}/versions/latest".
func GetSecret(ctx context.Context, location string) (string, error) {
	logger := httplog.LogEntry(ctx)
	logger.Info("Getting secret from secret manager", "location", location)
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
	logger.Info("Got DB password from secret manager successfully")
	// The secret payload is a byte array, so convert it to a string.
	return string(secret.Payload.Data), nil
}

func createURLTableIfNotExists(ctx context.Context, tx pgx.Tx, embeddingsTableName string) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946497); err != nil {
		return err
	}
	_, err := tx.Exec(ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS urls (
			url TEXT NOT NULL,
			collection_id uuid NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			document_id uuid,
			PRIMARY KEY (url, collection_id),
			FOREIGN KEY (collection_id) REFERENCES %s (uuid) ON DELETE CASCADE,
			FOREIGN KEY (document_id) REFERENCES %s (uuid) ON DELETE CASCADE)`,
			CollectionTableName, embeddingsTableName))
	if err != nil {
		return fmt.Errorf("failed to create urls table: %w", err)
	}
	return nil
}

func connect(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	// Configure the driver to connect to the database
	connString := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable",
		cfg.DB.Name, cfg.DB.User, cfg.DB.Pass)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Create a new dialer with any options
	d, err := cloudsqlconn.NewDialer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create dialer: %w", err)
	}

	// Tell the driver to use the Cloud SQL Go Connector to create connections
	config.ConnConfig.DialFunc = func(ctx context.Context, _ string, instance string) (net.Conn, error) {
		return d.Dial(ctx, cfg.DB.Instance)
	}

	// Interact with the driver directly as you normally would
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}
	return pool, nil
}
