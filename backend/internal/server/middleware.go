package server

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/5pirit5eal/swim-gen/internal/models"
	"github.com/go-chi/httplog/v2"
)

// HTTP middleware for Supabase authentication
//
// This middleware extracts authentication information from incoming HTTP requests,
// verifies it using Supabase, and adds the authenticated user information to the
// request context for downstream handlers to use.
func (rs *RAGService) SupabaseAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := httplog.LogEntry(r.Context())
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Debug("No user token provided")
			fallbackCtx := context.WithValue(r.Context(), models.UserIdCtxKey, "")
			next.ServeHTTP(w, r.WithContext(fallbackCtx))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Debug("Invalid authorization header format", "header", authHeader)
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		user, err := rs.auth.Auth.WithToken(token).GetUser()
		if err != nil {
			logger.Error("Failed to verify token", httplog.ErrAttr(err), "token", token)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userId := user.ID

		logger.Debug("Successfully verified user", "user_id", userId)
		httplog.LogEntrySetField(r.Context(), "user_id", slog.StringValue(userId.String()))

		ctx := context.WithValue(r.Context(), models.UserIdCtxKey, userId.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
