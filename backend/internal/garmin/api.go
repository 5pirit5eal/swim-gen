package garmin

// Garmin Connect API client for workout management, ported from:
// https://github.com/matin/garth/blob/main/src/garth/http.py
// https://github.com/cyberjunky/python-garminconnect/blob/master/garminconnect/workout.py

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultDomain  = "garmin.com"
	defaultTimeout = 10 * time.Second
)

// Client is an authenticated Garmin Connect API client.
type Client struct {
	httpClient  *http.Client
	OAuth1Token *OAuth1Token
	OAuth2Token *OAuth2Token
	Domain      string
	Timeout     time.Duration
	baseURL     string // override for testing; when set, used instead of https://connectapi.<domain>
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithDomain sets the Garmin domain (e.g. "garmin.com" or "garmin.cn").
func WithDomain(domain string) ClientOption {
	return func(c *Client) {
		c.Domain = domain
	}
}

// WithTimeout sets the HTTP request timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

// NewClient creates an authenticated Client by performing the SSO login flow.
func NewClient(email, password string, opts ...ClientOption) (*Client, error) {
	c := &Client{
		Domain:  defaultDomain,
		Timeout: defaultTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}

	oauth1, oauth2, err := Login(email, password, c.Domain)
	if err != nil {
		return nil, fmt.Errorf("garmin login: %w", err)
	}

	c.OAuth1Token = oauth1
	c.OAuth2Token = oauth2
	c.httpClient = &http.Client{Timeout: c.Timeout}
	return c, nil
}

// NewClientFromTokens creates a Client from existing OAuth tokens without logging in.
func NewClientFromTokens(oauth1 *OAuth1Token, oauth2 *OAuth2Token, opts ...ClientOption) *Client {
	c := &Client{
		Domain:      defaultDomain,
		Timeout:     defaultTimeout,
		OAuth1Token: oauth1,
		OAuth2Token: oauth2,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.httpClient = &http.Client{Timeout: c.Timeout}
	return c
}

// RefreshOAuth2 exchanges the OAuth1 token for a fresh OAuth2 token.
func (c *Client) RefreshOAuth2() error {
	if c.OAuth1Token == nil {
		return fmt.Errorf("OAuth1 token is required for OAuth2 refresh")
	}
	oauth2, err := Exchange(c.OAuth1Token, c.Domain, c.httpClient.Transport)
	if err != nil {
		return fmt.Errorf("refreshing OAuth2 token: %w", err)
	}
	c.OAuth2Token = oauth2
	return nil
}

// connectAPI performs an authenticated request to the Garmin Connect API.
// It automatically refreshes the OAuth2 token if expired.
func (c *Client) connectAPI(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if c.OAuth2Token == nil {
		return nil, fmt.Errorf("not authenticated: OAuth2 token is nil")
	}

	// Auto-refresh expired OAuth2 token
	if c.OAuth2Token.Expired() {
		if err := c.RefreshOAuth2(); err != nil {
			return nil, fmt.Errorf("auto-refresh failed: %w", err)
		}
	}

	var apiURL string
	if c.baseURL != "" {
		apiURL = c.baseURL + path
	} else {
		apiURL = fmt.Sprintf("https://connectapi.%s%s", c.Domain, path)
	}
	req, err := http.NewRequestWithContext(ctx, method, apiURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating API request: %w", err)
	}

	req.Header.Set("Authorization", c.OAuth2Token.String())
	req.Header.Set("User-Agent", userAgent)

	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request to %s: %w", path, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request to %s failed (status %d): %s", path, resp.StatusCode, string(respBody))
	}

	return resp, nil
}

// UploadWorkout uploads a workout JSON payload to Garmin Connect.
// The workoutJSON should be a valid Garmin workout JSON object or array.
func (c *Client) UploadWorkout(ctx context.Context, workoutJSON []byte) (map[string]any, error) {
	resp, err := c.connectAPI(ctx, http.MethodPost, "/workout-service/workout", bytes.NewReader(workoutJSON))
	if err != nil {
		return nil, fmt.Errorf("uploading workout: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding workout upload response: %w", err)
	}
	return result, nil
}

// GetWorkouts returns workouts starting at offset start with at most limit results.
func (c *Client) GetWorkouts(ctx context.Context, start, limit int) ([]map[string]any, error) {
	if start < 0 {
		return nil, fmt.Errorf("start must be non-negative")
	}
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive")
	}

	path := fmt.Sprintf("/workout-service/workouts?start=%d&limit=%d", start, limit)
	resp, err := c.connectAPI(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting workouts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var result []map[string]any
	// Depending on Garmin's exact response structure, you might need to adjust this from []map[string]any to just map[string]any if they wrap the array.
	// Commonly they return a direct JSON array for workouts.
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding workouts response: %w", err)
	}
	return result, nil
}

// GetWorkoutByID returns a workout by its ID.
func (c *Client) GetWorkoutByID(ctx context.Context, workoutID int64) (map[string]any, error) {
	if workoutID <= 0 {
		return nil, fmt.Errorf("workoutID must be positive")
	}

	path := fmt.Sprintf("/workout-service/workout/%d", workoutID)
	resp, err := c.connectAPI(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("getting workout %d: %w", workoutID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding workout response: %w", err)
	}
	return result, nil
}

// SaveTokens persists the client's tokens to the given directory.
// TODO: Implement database storage for tokens to support multi-instance deployments.
// The database implementation should replace the file-based storage when ready.
func (c *Client) SaveTokens(dir string) error {
	return SaveTokens(dir, c.OAuth1Token, c.OAuth2Token)
}

// LoadTokens reads tokens from the given directory and configures the client.
// TODO: Implement database storage for tokens to support multi-instance deployments.
// The database implementation should replace the file-based storage when ready.
func (c *Client) LoadTokens(dir string) error {
	oauth1, oauth2, err := LoadTokens(dir)
	if err != nil {
		return err
	}
	c.OAuth1Token = oauth1
	c.OAuth2Token = oauth2
	return nil
}

// SetHTTPClient replaces the underlying HTTP client (useful for testing).
func (c *Client) SetHTTPClient(hc *http.Client) {
	c.httpClient = hc
}

// SetBaseURL overrides the API base URL (useful for testing with httptest).
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}
