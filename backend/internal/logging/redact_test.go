package logging

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"testing"
)

func TestRedactorRemovesSensitiveValues(t *testing.T) {
	secret := "service-role-secret"
	redactor := NewRedactor(secret)

	var output bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: redactor.ReplaceAttr,
	}))
	logger.Debug("diagnostic payload", "message", "Bearer bearer-secret", "password", "database-password", "secret", secret, "err", errors.New("token=error-token"), slog.Group("nested", slog.String("cookie", "cookie-value")))

	result := output.String()
	for _, value := range []string{secret, "bearer-secret", "database-password", "error-token", "cookie-value"} {
		if strings.Contains(result, value) {
			t.Fatalf("log contains sensitive value %q: %s", value, result)
		}
	}
	for _, value := range []string{"[REDACTED]", "diagnostic payload"} {
		if !strings.Contains(result, value) {
			t.Fatalf("log does not contain %q: %s", value, result)
		}
	}
}

func TestRedactorKeepsPayloadsAtDebugButNotInfo(t *testing.T) {
	var debugOutput bytes.Buffer
	debugLogger := NewTextLogger(&debugOutput, slog.LevelDebug)
	debugLogger.Debug("chat payload", "message", "make this harder")
	if !strings.Contains(debugOutput.String(), "make this harder") {
		t.Fatalf("debug log did not contain payload: %s", debugOutput.String())
	}

	var infoOutput bytes.Buffer
	infoLogger := NewTextLogger(&infoOutput, slog.LevelInfo)
	infoLogger.Debug("chat payload", "message", "make this harder")
	if strings.Contains(infoOutput.String(), "make this harder") {
		t.Fatalf("info log contained debug payload: %s", infoOutput.String())
	}
}

func TestRedactorRemovesBearerAndJWTFromMessages(t *testing.T) {
	redactor := NewRedactor()
	value := redactor.RedactString("Authorization: Bearer bearer-token and jwt eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyIn0.signature123")

	if strings.Contains(value, "bearer-token") || strings.Contains(value, "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ1c2VyIn0.signature123") {
		t.Fatalf("sensitive message content was not redacted: %s", value)
	}
}

func TestRedactorRemovesSignedURLs(t *testing.T) {
	redactor := NewRedactor()
	value := redactor.RedactString("https://storage.example.test/plan.pdf?X-Goog-Algorithm=GOOG4-RSA-SHA256&X-Goog-Signature=signature-value")

	if strings.Contains(value, "signature-value") || strings.Contains(value, "storage.example.test") {
		t.Fatalf("signed URL was not redacted: %s", value)
	}
}
