package rag

import (
	"context"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/5pirit5eal/swim-rag/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

const (
	CollectionTableName string = "documents"
)

type RAGDB struct {
	Conn  pgvector.PGXConn
	Store *pgvector.Store
}

func NewStore(client embeddings.EmbedderClient, ctx context.Context, cfg models.Config) (*RAGDB, error) {
	// Load the database password from Google Secret Manager
	if cfg.DB.Pass == "" {
		pass, err := GetSecret(ctx, cfg.DB.PassLocation)
		if err != nil {
			return nil, err
		}
		cfg.DB.Pass = pass
	}
	log.Println("Got DB password successfully")

	// Create an embedder
	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}

	log.Println("Creating database connection...")
	// Initialize the database connection
	// TODO: connect via cloud sql proxy (Unix socket or TCP) or directly via URL
	cfg.DB.Instance = fmt.Sprintf("host=/tmp/cloudsql/%s database=%s user=%s password=%s",
		cfg.DB.Instance, cfg.DB.Name, cfg.DB.User, cfg.DB.Pass)
	conn, err := pgxpool.New(ctx, cfg.DB.Instance)
	if err != nil {
		return nil, err
	}
	log.Println("Database connection created successfully")

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
	log.Println("Created store successfully")
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
	return &RAGDB{Store: &store, Conn: conn}, nil
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
