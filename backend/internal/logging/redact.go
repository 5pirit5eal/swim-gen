package logging

import (
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
)

const redactedValue = "[REDACTED]"

var (
	bearerPattern    = regexp.MustCompile(`(?i)\bBearer\s+[^\s,;]+`)
	jwtPattern       = regexp.MustCompile(`\b[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`)
	signedURLPattern = regexp.MustCompile(`(?i)https?://[^\s"'<>]*(?:x-goog-signature|x-goog-credential|googleaccessid|x-amz-signature|x-amz-credential|signature|token)=[^\s&"'<>]+[^\s"'<>]*`)
	keyValuePattern  = regexp.MustCompile(`(?i)\b(password|passwd|pwd|secret|token|api[-_]?key|authorization|cookie|signature|signed[-_]?url)\b(\s*[:=]\s*)("[^"]*"|'[^']*'|[^\s,;}]+)`)
)

// Redactor removes credentials from structured log attributes and messages.
// Payload logs remain available at DEBUG, but embedded credentials are still removed.
type Redactor struct {
	secrets []string
}

func NewRedactor(secrets ...string) Redactor {
	filtered := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		if secret != "" {
			filtered = append(filtered, secret)
		}
	}
	return Redactor{secrets: filtered}
}

// ReplaceAttr is compatible with slog.HandlerOptions.ReplaceAttr and
// httplog.Options.ReplaceAttrsOverride.
func (r Redactor) ReplaceAttr(_ []string, attr slog.Attr) slog.Attr {
	if isSensitiveKey(attr.Key) {
		return slog.String(attr.Key, redactedValue)
	}

	if attr.Value.Kind() == slog.KindGroup {
		group := attr.Value.Group()
		for i := range group {
			group[i] = r.ReplaceAttr(nil, group[i])
		}
		return slog.Attr{Key: attr.Key, Value: slog.GroupValue(group...)}
	}

	if attr.Key == "err" || attr.Key == "error" {
		if err, ok := attr.Value.Any().(error); ok {
			return slog.String(attr.Key, r.RedactString(err.Error()))
		}
	}

	switch attr.Value.Kind() {
	case slog.KindString:
		return slog.String(attr.Key, r.RedactString(attr.Value.String()))
	case slog.KindAny:
		if value := attr.Value.Any(); value != nil {
			return slog.Any(attr.Key, r.RedactString(fmt.Sprint(value)))
		}
	}

	return attr
}

// RedactString removes known secrets and common bearer/credential patterns.
func (r Redactor) RedactString(value string) string {
	for _, secret := range r.secrets {
		value = strings.ReplaceAll(value, secret, redactedValue)
	}

	value = bearerPattern.ReplaceAllString(value, "Bearer "+redactedValue)
	value = jwtPattern.ReplaceAllString(value, redactedValue)
	value = signedURLPattern.ReplaceAllString(value, redactedValue)
	return keyValuePattern.ReplaceAllString(value, "${1}${2}"+redactedValue)
}

func isSensitiveKey(key string) bool {
	key = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(key, "-", "_"), " ", "_"))
	for _, part := range []string{
		"authorization", "cookie", "password", "passwd", "pwd", "secret",
		"token", "jwt", "api_key", "anon_key", "service_role", "signed_url",
		"signature",
	} {
		if strings.Contains(key, part) {
			return true
		}
	}
	return false
}

func NewTextLogger(writer io.Writer, level slog.Level, secrets ...string) *slog.Logger {
	redactor := NewRedactor(secrets...)
	return slog.New(slog.NewTextHandler(writer, &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: redactor.ReplaceAttr,
	}))
}
