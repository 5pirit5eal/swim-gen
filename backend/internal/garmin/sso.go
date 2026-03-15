package garmin

// Garmin SSO authentication flow, ported from:
// https://github.com/matin/garth/blob/main/src/garth/sso.py
// https://github.com/matin/garth/blob/main/src/garth/auth_tokens.py

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	csrfRE  = regexp.MustCompile(`name="_csrf"\s+value="(.+?)"`)
	titleRE = regexp.MustCompile(`<title>(.+?)</title>`)
)

const (
	oauthConsumerURL = "https://thegarth.s3.amazonaws.com/oauth_consumer.json"
	userAgent        = "com.garmin.android.apps.connectmobile"
)

// OAuth1Token holds the OAuth 1.0 credentials obtained from Garmin SSO.
type OAuth1Token struct {
	OAuthToken       string `json:"oauth_token"`
	OAuthTokenSecret string `json:"oauth_token_secret"`
	MFAToken         string `json:"mfa_token,omitempty"`
	Domain           string `json:"domain,omitempty"`
}

// OAuth2Token holds the OAuth 2.0 credentials exchanged from an OAuth1Token.
type OAuth2Token struct {
	Scope                  string `json:"scope"`
	JTI                    string `json:"jti"`
	TokenType              string `json:"token_type"`
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	ExpiresIn              int64  `json:"expires_in"`
	ExpiresAt              int64  `json:"expires_at"`
	RefreshTokenExpiresIn  int64  `json:"refresh_token_expires_in"`
	RefreshTokenExpiresAt  int64  `json:"refresh_token_expires_at"`
}

// Expired reports whether the access token has expired.
func (t *OAuth2Token) Expired() bool {
	return t.ExpiresAt < time.Now().Unix()
}

// RefreshExpired reports whether the refresh token has expired.
func (t *OAuth2Token) RefreshExpired() bool {
	return t.RefreshTokenExpiresAt < time.Now().Unix()
}

// String returns the Authorization header value: "Bearer <access_token>".
func (t *OAuth2Token) String() string {
	return fmt.Sprintf("Bearer %s", t.AccessToken)
}

// oauthConsumer holds the consumer key and secret fetched from Garmin's S3.
type oauthConsumer struct {
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
}

// fetchOAuthConsumer retrieves the OAuth consumer key and secret from Garmin.
func fetchOAuthConsumer(httpClient *http.Client) (*oauthConsumer, error) {
	resp, err := httpClient.Get(oauthConsumerURL)
	if err != nil {
		return nil, fmt.Errorf("fetching OAuth consumer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetching OAuth consumer: status %d", resp.StatusCode)
	}

	var consumer oauthConsumer
	if err := json.NewDecoder(resp.Body).Decode(&consumer); err != nil {
		return nil, fmt.Errorf("decoding OAuth consumer: %w", err)
	}
	return &consumer, nil
}

// Login performs the Garmin SSO login flow and returns OAuth1 and OAuth2 tokens.
// The domain is typically "garmin.com".
func Login(email, password, domain string) (*OAuth1Token, *OAuth2Token, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating cookie jar: %w", err)
	}
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	ssoBase := fmt.Sprintf("https://sso.%s/sso", domain)
	ssoEmbed := ssoBase + "/embed"

	embedParams := url.Values{
		"id":          {"gauth-widget"},
		"embedWidget": {"true"},
		"gauthHost":   {ssoBase},
	}

	signinParams := url.Values{
		"id":                              {"gauth-widget"},
		"embedWidget":                     {"true"},
		"gauthHost":                       {ssoEmbed},
		"service":                         {ssoEmbed},
		"source":                          {ssoEmbed},
		"redirectAfterAccountLoginUrl":    {ssoEmbed},
		"redirectAfterAccountCreationUrl": {ssoEmbed},
	}

	// Step 1: Set cookies
	embedURL := fmt.Sprintf("%s/embed?%s", ssoBase, embedParams.Encode())
	resp, err := httpClient.Get(embedURL)
	if err != nil {
		return nil, nil, fmt.Errorf("setting SSO cookies: %w", err)
	}
	resp.Body.Close()

	// Step 2: Get CSRF token
	signinURL := fmt.Sprintf("%s/signin?%s", ssoBase, signinParams.Encode())
	req, err := http.NewRequest(http.MethodGet, signinURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating signin request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching signin page: %w", err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("reading signin page: %w", err)
	}

	csrfToken, err := GetCSRFToken(string(body))
	if err != nil {
		return nil, nil, err
	}

	// Step 3: Submit login form
	formData := url.Values{
		"username": {email},
		"password": {password},
		"embed":    {"true"},
		"_csrf":    {csrfToken},
	}
	req, err = http.NewRequest(
		http.MethodPost,
		signinURL,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("submitting login: %w", err)
	}
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("reading login response: %w", err)
	}

	title, err := GetTitle(string(body))
	if err != nil {
		return nil, nil, err
	}
	if title != "Success" {
		return nil, nil, fmt.Errorf("login failed: unexpected title %q", title)
	}

	// Step 4: Extract ticket
	ticketRE := regexp.MustCompile(`embed\?ticket=([^"]+)"`)
	matches := ticketRE.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return nil, nil, fmt.Errorf("could not find ticket in login response")
	}
	ticket := matches[1]

	// Step 5: Get OAuth1 token
	oauth1, err := GetOAuth1Token(ticket, domain, httpClient.Transport)
	if err != nil {
		return nil, nil, fmt.Errorf("getting OAuth1 token: %w", err)
	}

	// Step 6: Exchange for OAuth2 token
	oauth2, err := Exchange(oauth1, domain, httpClient.Transport)
	if err != nil {
		return nil, nil, fmt.Errorf("exchanging for OAuth2 token: %w", err)
	}

	return oauth1, oauth2, nil
}

// GetOAuth1Token obtains an OAuth1 token using the SSO ticket.
func GetOAuth1Token(ticket, domain string, transport http.RoundTripper) (*OAuth1Token, error) {
	consumer, err := fetchOAuthConsumer(&http.Client{Transport: transport, Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}

	loginURL := fmt.Sprintf("https://sso.%s/sso/embed", domain)
	baseURL := fmt.Sprintf("https://connectapi.%s/oauth-service/oauth/preauthorized", domain)

	params := url.Values{
		"ticket":             {ticket},
		"login-url":          {loginURL},
		"accepts-mfa-tokens": {"true"},
	}

	fullURL := baseURL + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating OAuth1 request: %w", err)
	}

	// Sign the request with OAuth1 (consumer-only, no token yet)
	signOAuth1Request(req, consumer, "", "")
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OAuth1 preauthorized request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OAuth1 preauthorized failed (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading OAuth1 response: %w", err)
	}

	parsed, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("parsing OAuth1 response: %w", err)
	}

	return &OAuth1Token{
		OAuthToken:       parsed.Get("oauth_token"),
		OAuthTokenSecret: parsed.Get("oauth_token_secret"),
		MFAToken:         parsed.Get("mfa_token"),
		Domain:           domain,
	}, nil
}

// Exchange exchanges an OAuth1 token for an OAuth2 token.
func Exchange(oauth1 *OAuth1Token, domain string, transport http.RoundTripper) (*OAuth2Token, error) {
	consumer, err := fetchOAuthConsumer(&http.Client{Transport: transport, Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}

	baseURL := fmt.Sprintf("https://connectapi.%s/oauth-service/oauth/exchange/user/2.0", domain)

	formData := url.Values{}
	if oauth1.MFAToken != "" {
		formData.Set("mfa_token", oauth1.MFAToken)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating OAuth2 exchange request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)

	// Sign with OAuth1 (consumer + token)
	signOAuth1Request(req, consumer, oauth1.OAuthToken, oauth1.OAuthTokenSecret)

	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OAuth2 exchange request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OAuth2 exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding OAuth2 response: %w", err)
	}

	now := time.Now().Unix()
	expiresIn := jsonInt64(tokenResp, "expires_in")
	refreshExpiresIn := jsonInt64(tokenResp, "refresh_token_expires_in")

	return &OAuth2Token{
		Scope:                  jsonString(tokenResp, "scope"),
		JTI:                    jsonString(tokenResp, "jti"),
		TokenType:              jsonString(tokenResp, "token_type"),
		AccessToken:            jsonString(tokenResp, "access_token"),
		RefreshToken:           jsonString(tokenResp, "refresh_token"),
		ExpiresIn:              expiresIn,
		ExpiresAt:              now + expiresIn,
		RefreshTokenExpiresIn:  refreshExpiresIn,
		RefreshTokenExpiresAt:  now + refreshExpiresIn,
	}, nil
}

// SaveTokens persists the OAuth1 and OAuth2 tokens to the given directory as JSON files.
// TODO: Implement database storage for tokens to support multi-instance deployments.
// The database implementation should replace the file-based storage when ready.
func SaveTokens(dir string, oauth1 *OAuth1Token, oauth2 *OAuth2Token) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating token directory: %w", err)
	}

	if err := writeJSON(filepath.Join(dir, "oauth1_token.json"), oauth1); err != nil {
		return fmt.Errorf("saving OAuth1 token: %w", err)
	}
	if err := writeJSON(filepath.Join(dir, "oauth2_token.json"), oauth2); err != nil {
		return fmt.Errorf("saving OAuth2 token: %w", err)
	}
	return nil
}

// LoadTokens reads OAuth1 and OAuth2 tokens from the given directory.
// TODO: Implement database storage for tokens to support multi-instance deployments.
// The database implementation should replace the file-based storage when ready.
func LoadTokens(dir string) (*OAuth1Token, *OAuth2Token, error) {
	var oauth1 OAuth1Token
	if err := readJSON(filepath.Join(dir, "oauth1_token.json"), &oauth1); err != nil {
		return nil, nil, fmt.Errorf("loading OAuth1 token: %w", err)
	}

	var oauth2 OAuth2Token
	if err := readJSON(filepath.Join(dir, "oauth2_token.json"), &oauth2); err != nil {
		return nil, nil, fmt.Errorf("loading OAuth2 token: %w", err)
	}

	return &oauth1, &oauth2, nil
}

// --- HTML parsing helpers ---

// GetCSRFToken extracts the CSRF token from the SSO HTML page.
func GetCSRFToken(html string) (string, error) {
	matches := csrfRE.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find CSRF token in response")
	}
	return matches[1], nil
}

// GetTitle extracts the <title> from an HTML page.
func GetTitle(html string) (string, error) {
	matches := titleRE.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find title in response")
	}
	return matches[1], nil
}

// --- OAuth1 signing ---

// signOAuth1Request signs an HTTP request using OAuth 1.0a HMAC-SHA1.
// If oauthToken and oauthTokenSecret are empty, only the consumer credentials are used.
func signOAuth1Request(req *http.Request, consumer *oauthConsumer, oauthToken, oauthTokenSecret string) {
	nonce := generateNonce()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	oauthParams := map[string]string{
		"oauth_consumer_key":     consumer.ConsumerKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        timestamp,
		"oauth_version":          "1.0",
	}

	if oauthToken != "" {
		oauthParams["oauth_token"] = oauthToken
	}

	// Collect all parameters (query + oauth) for signature base string
	allParams := url.Values{}
	for k, v := range oauthParams {
		allParams.Set(k, v)
	}
	for k, vs := range req.URL.Query() {
		for _, v := range vs {
			allParams.Add(k, v)
		}
	}

	// Build signature base string
	baseStringURL := strings.SplitN(req.URL.String(), "?", 2)[0]
	paramString := sortedParamString(allParams)
	baseString := strings.Join([]string{
		strings.ToUpper(req.Method),
		url.QueryEscape(baseStringURL),
		url.QueryEscape(paramString),
	}, "&")

	// Sign
	signingKey := url.QueryEscape(consumer.ConsumerSecret) + "&" + url.QueryEscape(oauthTokenSecret)
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(baseString))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	oauthParams["oauth_signature"] = signature

	// Build Authorization header
	var parts []string
	for k, v := range oauthParams {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, url.QueryEscape(k), url.QueryEscape(v)))
	}
	sort.Strings(parts)
	req.Header.Set("Authorization", "OAuth "+strings.Join(parts, ", "))
}

// sortedParamString builds the sorted parameter string for OAuth1 signature.
func sortedParamString(params url.Values) string {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		vs := params[k]
		sort.Strings(vs)
		for _, v := range vs {
			pairs = append(pairs, url.QueryEscape(k)+"="+url.QueryEscape(v))
		}
	}
	return strings.Join(pairs, "&")
}

// generateNonce creates a random nonce for OAuth1.
func generateNonce() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// --- JSON file helpers ---

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// --- JSON parsing helpers for untyped maps ---

func jsonString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func jsonInt64(m map[string]any, key string) int64 {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		}
	}
	return 0
}
