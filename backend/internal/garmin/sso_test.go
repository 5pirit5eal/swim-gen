package garmin_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2Token_Expired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		token := &garmin.OAuth2Token{
			ExpiresAt: time.Now().Unix() + 3600,
		}
		assert.False(t, token.Expired())
	})

	t.Run("expired", func(t *testing.T) {
		token := &garmin.OAuth2Token{
			ExpiresAt: time.Now().Unix() - 1,
		}
		assert.True(t, token.Expired())
	})

	t.Run("boundary - ExpiresAt equals now", func(t *testing.T) {
		token := &garmin.OAuth2Token{
			ExpiresAt: time.Now().Unix(),
		}
		// ExpiresAt == now is NOT expired (strict < comparison)
		assert.False(t, token.Expired())
	})
}

func TestOAuth2Token_RefreshExpired(t *testing.T) {
	t.Run("not expired", func(t *testing.T) {
		token := &garmin.OAuth2Token{
			RefreshTokenExpiresAt: time.Now().Unix() + 86400,
		}
		assert.False(t, token.RefreshExpired())
	})

	t.Run("expired", func(t *testing.T) {
		token := &garmin.OAuth2Token{
			RefreshTokenExpiresAt: time.Now().Unix() - 1,
		}
		assert.True(t, token.RefreshExpired())
	})
}

func TestOAuth2Token_String(t *testing.T) {
	token := &garmin.OAuth2Token{
		TokenType:   "Bearer",
		AccessToken: "test_access_token_123",
	}
	assert.Equal(t, "Bearer test_access_token_123", token.String())
}

func TestGetCSRFToken(t *testing.T) {
	t.Run("valid CSRF token", func(t *testing.T) {
		html := `<html><input name="_csrf" value="abc123def456"></html>`
		token, err := garmin.GetCSRFToken(html)
		require.NoError(t, err)
		assert.Equal(t, "abc123def456", token)
	})

	t.Run("missing CSRF token", func(t *testing.T) {
		html := `<html><body>no csrf here</body></html>`
		_, err := garmin.GetCSRFToken(html)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "CSRF")
	})

	t.Run("CSRF with extra whitespace", func(t *testing.T) {
		html := `<input name="_csrf"    value="xyz789">`
		token, err := garmin.GetCSRFToken(html)
		require.NoError(t, err)
		assert.Equal(t, "xyz789", token)
	})
}

func TestGetTitle(t *testing.T) {
	t.Run("success title", func(t *testing.T) {
		html := `<html><head><title>Success</title></head></html>`
		title, err := garmin.GetTitle(html)
		require.NoError(t, err)
		assert.Equal(t, "Success", title)
	})

	t.Run("MFA challenge title", func(t *testing.T) {
		html := `<html><head><title>GARMIN > MFA Challenge</title></head></html>`
		title, err := garmin.GetTitle(html)
		require.NoError(t, err)
		assert.Contains(t, title, "MFA")
	})

	t.Run("missing title", func(t *testing.T) {
		html := `<html><body>no title</body></html>`
		_, err := garmin.GetTitle(html)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "title")
	})
}

func TestSaveAndLoadTokens(t *testing.T) {
	dir := t.TempDir()

	oauth1 := &garmin.OAuth1Token{
		OAuthToken:       "test_token",
		OAuthTokenSecret: "test_secret",
		MFAToken:         "",
		Domain:           "garmin.com",
	}
	oauth2 := &garmin.OAuth2Token{
		Scope:                  "CONNECT_READ CONNECT_WRITE",
		JTI:                    "test_jti",
		TokenType:              "Bearer",
		AccessToken:            "test_access",
		RefreshToken:           "test_refresh",
		ExpiresIn:              3600,
		ExpiresAt:              time.Now().Unix() + 3600,
		RefreshTokenExpiresIn:  7776000,
		RefreshTokenExpiresAt:  time.Now().Unix() + 7776000,
	}

	err := garmin.SaveTokens(dir, oauth1, oauth2)
	require.NoError(t, err)

	// Verify files exist
	assert.FileExists(t, filepath.Join(dir, "oauth1_token.json"))
	assert.FileExists(t, filepath.Join(dir, "oauth2_token.json"))

	// Load tokens
	loaded1, loaded2, err := garmin.LoadTokens(dir)
	require.NoError(t, err)

	assert.Equal(t, oauth1.OAuthToken, loaded1.OAuthToken)
	assert.Equal(t, oauth1.OAuthTokenSecret, loaded1.OAuthTokenSecret)
	assert.Equal(t, oauth1.Domain, loaded1.Domain)

	assert.Equal(t, oauth2.AccessToken, loaded2.AccessToken)
	assert.Equal(t, oauth2.RefreshToken, loaded2.RefreshToken)
	assert.Equal(t, oauth2.Scope, loaded2.Scope)
	assert.Equal(t, oauth2.ExpiresAt, loaded2.ExpiresAt)
}

func TestLoadTokens_MissingFiles(t *testing.T) {
	dir := t.TempDir()
	_, _, err := garmin.LoadTokens(dir)
	assert.Error(t, err)
}

func TestSaveTokens_CreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "subdir", "tokens")

	oauth1 := &garmin.OAuth1Token{OAuthToken: "t", OAuthTokenSecret: "s"}
	oauth2 := &garmin.OAuth2Token{AccessToken: "a"}

	err := garmin.SaveTokens(dir, oauth1, oauth2)
	require.NoError(t, err)

	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}
