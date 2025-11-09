package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/5pirit5eal/swim-gen/docs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/server"
)

// Package main provides the swim-gen API server
//
//	@title			Swim Gen API
//	@version		1.0
//	@description	A REST API for swim training plan management with RAG capabilities
//
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@BasePath	/
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT.
func main() {
	// command line flags
	envFile := flag.String("env", ".env", "path to .env file")
	flag.Parse()

	// Configure log to write to stdout
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := config.LoadConfig(filepath.Join(projectRoot, *envFile), true)
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

	router := setupRouter("/", ragServer, cfg, logger)

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
	logger := httplog.NewLogger("swim-gen", httplog.Options{
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

// Basic health check endpoint - moved to server.HealthHandler for comprehensive checks
// This is kept for backward compatibility with simple health checks
// @Summary Basic health check
// @Description Returns a simple OK response for basic health monitoring
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Router /health-basic [get]
func basicHealthHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		httplog.LogEntry(r.Context()).Error("Failed to write health check response", httplog.ErrAttr(err))
	}
}

// Setup of routes for the RAG service
func setupRouter(basePath string, ragServer *server.RAGService, cfg config.Config, logger *httplog.Logger) chi.Router {

	// Service
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger, []string{"/health"}))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route(basePath, func(r chi.Router) {
		r.Use(ragServer.SupabaseAuthMiddleware)
		r.Get("/health", ragServer.HealthHandler)
		r.Get("/health-basic", basicHealthHandler)
		r.Post("/add", ragServer.DonatePlanHandler)
		r.Post("/prompt", ragServer.GeneratePromptHandler)
		r.Post("/query", ragServer.QueryHandler)
		r.Get("/scrape", ragServer.ScrapeHandler)
		r.Post("/export-pdf", ragServer.PlanToPDFHandler)
		r.Post("/upsert-plan", ragServer.UpsertPlanHandler)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("0.0.0.0:"+cmp.Or(cfg.Port, "8080")+basePath+"swagger/doc.json"),
			httpSwagger.DeepLinking(true)),
		)
	})

	return r
}
