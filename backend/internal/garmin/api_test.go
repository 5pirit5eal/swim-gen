package garmin_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/5pirit5eal/swim-gen/internal/garmin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientFromTokens(t *testing.T) {
	oauth1 := &garmin.OAuth1Token{
		OAuthToken:       "test_token",
		OAuthTokenSecret: "test_secret",
		Domain:           "garmin.com",
	}
	oauth2 := &garmin.OAuth2Token{
		TokenType:   "Bearer",
		AccessToken: "test_access",
		ExpiresAt:   time.Now().Unix() + 3600,
	}

	client := garmin.NewClientFromTokens(oauth1, oauth2)
	assert.NotNil(t, client)
	assert.Equal(t, "garmin.com", client.Domain)
	assert.Equal(t, oauth1, client.OAuth1Token)
	assert.Equal(t, oauth2, client.OAuth2Token)
}

func TestNewClientFromTokens_WithOptions(t *testing.T) {
	oauth1 := &garmin.OAuth1Token{OAuthToken: "t", OAuthTokenSecret: "s"}
	oauth2 := &garmin.OAuth2Token{AccessToken: "a", ExpiresAt: time.Now().Unix() + 3600}

	client := garmin.NewClientFromTokens(
		oauth1, oauth2,
		garmin.WithDomain("garmin.cn"),
		garmin.WithTimeout(30*time.Second),
	)
	assert.Equal(t, "garmin.cn", client.Domain)
	assert.Equal(t, 30*time.Second, client.Timeout)
}

func TestUploadWorkout(t *testing.T) {
	workoutResponse := map[string]any{
		"workoutId":   float64(12345),
		"workoutName": "Test Swim Workout",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/workout-service/workout", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(workoutResponse)
	}))
	defer server.Close()

	// Create a client that points to the test server
	oauth1 := &garmin.OAuth1Token{OAuthToken: "t", OAuthTokenSecret: "s"}
	oauth2 := &garmin.OAuth2Token{
		TokenType:   "Bearer",
		AccessToken: "test_access_token",
		ExpiresAt:   time.Now().Unix() + 3600,
	}

	// Use a custom domain that resolves to the test server
	client := garmin.NewClientFromTokens(oauth1, oauth2)
	// Override httpClient to point to test server via custom transport
	client.SetHTTPClient(server.Client())
	client.SetBaseURL(server.URL)

	workoutJSON := []byte(`{"workoutName":"Test Swim Workout"}`)
	result, err := client.UploadWorkout(context.Background(), workoutJSON)
	require.NoError(t, err)
	assert.Equal(t, float64(12345), result["workoutId"])
	assert.Equal(t, "Test Swim Workout", result["workoutName"])
}

func TestClientSaveAndLoadTokens(t *testing.T) {
	dir := t.TempDir()

	oauth1 := &garmin.OAuth1Token{
		OAuthToken:       "token",
		OAuthTokenSecret: "secret",
		Domain:           "garmin.com",
	}
	oauth2 := &garmin.OAuth2Token{
		AccessToken:  "access",
		RefreshToken: "refresh",
		ExpiresAt:    time.Now().Unix() + 3600,
	}

	client := garmin.NewClientFromTokens(oauth1, oauth2)
	err := client.SaveTokens(dir)
	require.NoError(t, err)

	// Create a new client and load tokens
	newClient := garmin.NewClientFromTokens(
		&garmin.OAuth1Token{},
		&garmin.OAuth2Token{},
	)
	err = newClient.LoadTokens(dir)
	require.NoError(t, err)

	assert.Equal(t, oauth1.OAuthToken, newClient.OAuth1Token.OAuthToken)
	assert.Equal(t, oauth2.AccessToken, newClient.OAuth2Token.AccessToken)
}

func TestConnectAPI_RefreshesExpiredToken(t *testing.T) {
	var refreshCalled bool

	mockRoundTripper := func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "thegarth.s3.amazonaws.com" {
			// Mock oauth_consumer.json
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"consumer_key":"k","consumer_secret":"s"}`)),
			}, nil
		}
		if strings.Contains(req.URL.Path, "exchange/user/2.0") {
			refreshCalled = true
			tokenResp := `{"access_token":"new_token","expires_in":3600}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(tokenResp)),
			}, nil
		}
		if strings.HasSuffix(req.URL.Path, "/workout-service/workout") {
			// Validate that the request uses the *new* access token
			assert.Equal(t, "Bearer new_token", req.Header.Get("Authorization"))
			return &http.Response{StatusCode: http.StatusNoContent, Body: io.NopCloser(strings.NewReader(""))}, nil
		}
		return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader(""))}, nil
	}

	client := garmin.NewClientFromTokens(
		&garmin.OAuth1Token{OAuthToken: "t", OAuthTokenSecret: "s"},
		&garmin.OAuth2Token{AccessToken: "expired_token", ExpiresAt: time.Now().Unix() - 3600}, // manually expired
	)

	// Inject the mock transport
	customClient := &http.Client{
		Transport: &mockTransport{roundTripFunc: mockRoundTripper},
	}
	client.SetHTTPClient(customClient)

	// The client should detect the token is expired and refresh it automatically during the upload request
	_, err := client.UploadWorkout(context.Background(), []byte(`{}`))
	require.NoError(t, err)
	assert.True(t, refreshCalled, "expected token refresh to be triggered")
}

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

