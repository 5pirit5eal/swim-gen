package main

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/5pirit5eal/swim-gen/internal/config"
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
