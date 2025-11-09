package rag

import (
	"context"
	"fmt"
	"log/slog"

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
		pgvector.WithEmbeddingTableName(cfg.Embedding.Name), // Ensure this matches your SQL table name
		pgvector.WithCollectionTableName(CollectionTableName),
		// Separate the collections and documents by embedding model name
		pgvector.WithCollectionName(cfg.Embedding.Model),
		pgvector.WithEmbedder(embedder),
		pgvector.WithVectorDimensions(cfg.Embedding.Size),
	)
	if err != nil {
		return nil, err
	}
	slog.Info("Created langchaingo pgvector datastore successfully")
	// Check that the tables exist
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := checkScrapedTableExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := checkPlanTableExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := checkFeedbackTableExists(ctx, tx); err != nil {
		return nil, err
	}
	if err := checkDonatedPlanTableExists(ctx, tx); err != nil {
		return nil, err
	}
	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	slog.Info("All required tables exist")
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

func checkScrapedTableExists(ctx context.Context, tx pgx.Tx) error {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = $1
		)`, ScrapedTableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if scraped table exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("scraped table '%s' does not exist, please run migrations", ScrapedTableName)
	}
	return nil
}

func checkPlanTableExists(ctx context.Context, tx pgx.Tx) error {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = $1
		)`, PlanTableName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if plan table exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("plan table '%s' does not exist, please run migrations", PlanTableName)
	}
	return nil
}

func checkFeedbackTableExists(ctx context.Context, tx pgx.Tx) error {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = $1
		)`, FeedbackTable).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if feedback table exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("feedback table '%s' does not exist, please run migrations", FeedbackTable)
	}
	return nil
}

func checkDonatedPlanTableExists(ctx context.Context, tx pgx.Tx) error {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_name = $1
		)`, DonatedPlanTable).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if donated plan table exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("donated plan table '%s' does not exist, please run migrations", DonatedPlanTable)
	}
	return nil
}

func connect(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	// Configure the driver to connect to the database
	connString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s pool_max_conn_lifetime=30m",
		cfg.DB.Name, cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.SslMode)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}
	return pool, nil
}
