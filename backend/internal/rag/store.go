package rag

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/genai"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

const (
	CollectionTableName string = "embedders"
	PlanTableName       string = "plans"
	ScrapedTableName    string = "scraped"
	FeedbackTable       string = "feedback"
	DonatedPlanTable    string = "donations"
)

type RAGDB struct {
	Conn   pgvector.PGXConn
	Store  *pgvector.Store
	Client *genai.GoogleGenAIClient
	cfg    config.Config
}

func NewGoogleAIStore(ctx context.Context, cfg config.Config) (*RAGDB, error) {
	slog.Info("Initializing Google AI store with config", "cfg", slog.AnyValue(cfg))
	// Initialize the LLM client
	client, err := genai.NewGoogleGenAIClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	// Load the database password from Google Secret Manager
	if cfg.DB.Pass == "" {
		pass, err := GetSecret(ctx, cfg.DB.PassLocation)
		if err != nil {
			slog.Error("Failed to get DB password from secret manager", "error", err)
			return nil, err
		}
		cfg.DB.Pass = pass
	}
	slog.Info("Got DB password successfully")

	// Create an embedder
	embedder, err := embeddings.NewEmbedder(client)
	if err != nil {
		return nil, err
	}

	slog.Info("Creating database connection...")
	// Initialize the database connection
	conn, err := connect(ctx, cfg)
	if err != nil {
		return nil, err
	}
	slog.Info("Database connection created successfully")

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
	slog.Info("Created store successfully")
	// Create the URL table if it doesn't exist
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := createScrapedTableIfNotExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := createPlanTableIfNotExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := createFeedbackTableIfNotExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := createDonatedPlanTableIfNotExists(ctx, tx); err != nil {
		return nil, err
	}
	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	slog.Info("Setup URL table successfully")
	return &RAGDB{Store: &store, Conn: conn, Client: client, cfg: cfg}, nil
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
	slog.Info("Getting secret from secret manager", "location", location)
	// Create a new Secret Manager client
	// and access the secret version.
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = c.Close() }()
	secret, err := c.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: location,
	})
	if err != nil {
		return "", err
	}
	slog.Info("Got DB password from secret manager successfully")
	// The secret payload is a byte array, so convert it to a string.
	return string(secret.Payload.Data), nil
}

func createScrapedTableIfNotExists(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946497); err != nil {
		return err
	}
	_, err := tx.Exec(ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			url TEXT NOT NULL,
			collection_id uuid NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			plan_id uuid,
			PRIMARY KEY (url, collection_id),
			FOREIGN KEY (collection_id) REFERENCES %s (uuid) ON DELETE CASCADE)`,
			ScrapedTableName, CollectionTableName))
	if err != nil {
		return fmt.Errorf("failed to create scraped table: %w", err)
	}
	return nil
}

func createPlanTableIfNotExists(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946498); err != nil {
		return err
	}
	_, err := tx.Exec(ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			plan_id uuid NOT NULL DEFAULT gen_random_uuid(),
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			plan_table JSON NOT NULL,
			PRIMARY KEY (plan_id))`, PlanTableName))
	if err != nil {
		return fmt.Errorf("failed to create urls table: %w", err)
	}
	return nil
}

func createFeedbackTableIfNotExists(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946499); err != nil {
		return err
	}
	_, err := tx.Exec(ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			user_id uuid NOT NULL,
			plan_id uuid NOT NULL,
			rating INT NOT NULL,
			comment TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (user_id, plan_id))`, FeedbackTable))
	if err != nil {
		return fmt.Errorf("failed to create feedback table: %w", err)
	}
	return nil
}

func createDonatedPlanTableIfNotExists(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, "SELECT pg_advisory_xact_lock($1)", 1573678846307946500); err != nil {
		return err
	}
	_, err := tx.Exec(ctx,
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			user_id uuid NOT NULL,
			plan_id uuid NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			PRIMARY KEY (user_id, plan_id))`, DonatedPlanTable))
	if err != nil {
		return fmt.Errorf("failed to create donated table: %w", err)
	}
	return nil
}

func connect(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	// Configure the driver to connect to the database
	connString := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable pool_max_conns=50 pool_min_conns=5 pool_max_conn_lifetime=30m",
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
