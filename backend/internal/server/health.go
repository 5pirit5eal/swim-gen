package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
)

// HealthHandler performs comprehensive health checks including database and vector store connectivity
//
//	@Summary		Comprehensive health check
//	@Description	Returns the health status of the API including database and vector store connectivity
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	models.HealthStatus
//	@Failure		503	{object}	models.HealthStatus
//	@Router			/health [get]
func (s *RAGService) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	status := models.HealthStatus{
		Status:     "healthy",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Components: make(map[string]string),
	}

	// Check database connectivity
	if err := s.db.Conn.Ping(ctx); err != nil {
		httplog.LogEntry(ctx).Error("Database ping failed", httplog.ErrAttr(err))
		status.Components["database"] = "unhealthy"
		status.Status = "unhealthy"
	} else {
		status.Components["database"] = "healthy"

		// Try to get schema version
		var schemaVersion int
		err := s.db.Conn.QueryRow(ctx,
			"SELECT (value->>'version')::int FROM app_metadata WHERE key = 'schema_version'",
		).Scan(&schemaVersion)
		if err == nil {
			status.SchemaVersion = schemaVersion
		} else {
			httplog.LogEntry(ctx).Warn("Could not retrieve schema version", httplog.ErrAttr(err))
		}
	}

	// Check vector store readiness (collection exists)
	if s.db.Store != nil {
		// Try a simple similarity search with empty query to verify store is accessible
		ctx2, cancel2 := context.WithTimeout(ctx, 2*time.Second)
		defer cancel2()

		_, err := s.db.Store.SimilaritySearch(ctx2, "", 1)
		if err != nil {
			httplog.LogEntry(ctx).Warn("Vector store check failed", httplog.ErrAttr(err))
			status.Components["vector_store"] = "degraded"
			// Don't mark overall status as unhealthy - this is non-critical
		} else {
			status.Components["vector_store"] = "healthy"
		}
	} else {
		status.Components["vector_store"] = "not_initialized"
		httplog.LogEntry(ctx).Warn("Vector store is not initialized")
	}

	// Set HTTP status code
	statusCode := http.StatusOK
	if status.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(status); err != nil {
		httplog.LogEntry(ctx).Error("Failed to encode health response", httplog.ErrAttr(err))
	}
}
