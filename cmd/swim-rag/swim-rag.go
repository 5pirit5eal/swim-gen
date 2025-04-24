package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"

	"github.com/5pirit5eal/swim-rag/internal/config"
	"github.com/5pirit5eal/swim-rag/internal/scraper"
	"github.com/5pirit5eal/swim-rag/internal/server"
)

func main() {
	// Configure log to write to stdout
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.LoadConfig(filepath.Join(projectRoot, ".env"), true)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := setupLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ragServer, err := server.NewRAGService(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer ragServer.Close()
	scraper, err := scraper.NewScraper(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer scraper.Close()
	router := newRouter("/", ragServer, scraper, cfg, logger)

	port := cmp.Or(cfg.Port, "8080")
	address := "0.0.0.0:" + port
	logger.Info("Starting server", "listening on", address)
	if err := http.ListenAndServe(address, router); err != nil {
		logger.Error("Server stopped with error", httplog.ErrAttr(err))
		log.Fatal(err)
	}
}

func setupLogger(cfg config.Config) (*httplog.Logger, error) {
	// Check if we are in Cloud Run by looking for the K_SERVICE env var
	var j bool
	if _, exists := os.LookupEnv("K_SERVICE"); exists {
		j = true
	}
	levelMap := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}
	logger := httplog.NewLogger("swim-rag", httplog.Options{
		LogLevel: levelMap[cfg.LogLevel],
		JSON:     j,
		Concise:  true,
		// RequestHeaders:   true,
		// ResponseHeaders:  true,
		MessageFieldName: "message",
		LevelFieldName:   "severity",
		TimeFieldFormat:  time.RFC3339,
		QuietDownRoutes: []string{
			"/",
			"/health",
		},
		QuietDownPeriod: 10 * time.Second,
	})

	if logger == nil {
		return nil, fmt.Errorf("failed to create logger")
	}

	return logger, nil
}

// Setup of routes for the RAG service
func newRouter(basePath string, ragServer *server.RAGService, scraper *scraper.Scraper, cfg config.Config, logger *httplog.Logger) chi.Router {

	// Service
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger, []string{"/health"}))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route(basePath, func(r chi.Router) {
		r.Post("/add", ragServer.AddDocumentsHandler)
		r.Post("/query", ragServer.QueryHandler)
		r.Get("/scrape", scraper.ScrapeHandler)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
	})

	return r
}
