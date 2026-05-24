package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend/internal/usecase"
)

const accountsBaseURL = "https://accounts.spotify.com"

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
}

type OAuthToken struct {
	AccessToken  string
	TokenType    string
	Scope        string
	RefreshToken string
	ExpiresIn    int
	ExpiresAt    time.Time
}

type Profile struct {
	ID          string
	DisplayName string
}

type OAuthClient struct {
	config     OAuthConfig
	httpClient HTTPClient
}

func NewOAuthClient(config OAuthConfig, httpClient HTTPClient) *OAuthClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &OAuthClient{
		config:     config,
		httpClient: httpClient,
	}
}

func (c *OAuthClient) AuthorizeURL(state string) (string, error) {
	if strings.TrimSpace(state) == "" {
		return "", usecase.ErrInvalidInput
	}

	u, err := url.Parse(accountsBaseURL + "/authorize")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.config.ClientID)
	q.Set("redirect_uri", c.config.RedirectURI)
	q.Set("state", state)
	if len(c.config.Scopes) > 0 {
		q.Set("scope", strings.Join(c.config.Scopes, " "))
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (c *OAuthClient) ExchangeCode(ctx context.Context, code string) (*OAuthToken, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", c.config.RedirectURI)
	return c.tokenRequest(ctx, values)
}

func (c *OAuthClient) RefreshToken(ctx context.Context, refreshToken string) (*OAuthToken, error) {
	values := url.Values{}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", refreshToken)
	return c.tokenRequest(ctx, values)
}

func (c *OAuthClient) GetCurrentUserProfile(ctx context.Context, accessToken string) (*Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultBaseURL+"/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, usecase.ErrPermissionDenied
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("spotify profile error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &Profile{
		ID:          raw.ID,
		DisplayName: raw.DisplayName,
	}, nil
}

func (c *OAuthClient) tokenRequest(ctx context.Context, values url.Values) (*OAuthToken, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		accountsBaseURL+"/api/token",
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+basicAuth(c.config.ClientID, c.config.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusBadRequest {
		return nil, usecase.ErrPermissionDenied
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("spotify token error: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var token OAuthToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}
	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	return &token, nil
}

func basicAuth(clientID string, clientSecret string) string {
	return base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
}
