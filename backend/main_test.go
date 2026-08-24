package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/config"
	"github.com/5pirit5eal/swim-gen/internal/server"
)

func TestSetupLoggerDisablesHeadersAndRedactsConfiguredSecrets(t *testing.T) {
	cfg := config.Config{
		LogLevel: "DEBUG",
	}
	cfg.DB.Pass = "database-password"
	cfg.SB.AnonKey = "supabase-anon-key"
	cfg.SB.ServiceRoleKey = "supabase-service-role-key"

	logger, err := setupLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}

	if logger.Options.RequestHeaders || logger.Options.ResponseHeaders {
		t.Fatalf("request and response headers must not be logged")
	}

	attr := logger.Options.ReplaceAttrsOverride(nil, slog.String("diagnostic", "database-password"))
	if strings.Contains(attr.Value.String(), "database-password") {
		t.Fatalf("configured database password was not redacted")
	}
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := securityHeadersMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	request := httptest.NewRequest(http.MethodGet, "/api/private", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("expected API responses to be non-cacheable, got %q", response.Header().Get("Cache-Control"))
	}
	for name, expected := range map[string]string{
		"X-Content-Type-Options":     "nosniff",
		"X-Frame-Options":            "SAMEORIGIN",
		"Referrer-Policy":            "strict-origin-when-cross-origin",
		"Permissions-Policy":         "camera=(), microphone=(), geolocation=(), payment=()",
		"Cross-Origin-Opener-Policy": "same-origin",
	} {
		if actual := response.Header().Get(name); actual != expected {
			t.Errorf("expected %s=%q, got %q", name, expected, actual)
		}
	}
}

func TestSwaggerIsDisabledByDefault(t *testing.T) {
	cfg := config.Config{LogLevel: "INFO"}
	logger, err := setupLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}

	router := setupRouter("/", &server.RAGService{}, cfg, logger)
	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected Swagger to be disabled, got status %d", response.Code)
	}
}

func TestSwaggerCanBeEnabledForDevelopment(t *testing.T) {
	cfg := config.Config{LogLevel: "INFO", SwaggerEnabled: true}
	logger, err := setupLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}

	router := setupRouter("/", &server.RAGService{}, cfg, logger)
	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected enabled Swagger UI, got status %d", response.Code)
	}
}
