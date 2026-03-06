package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/boscod/responsewatch-cli/internal/config"
	"github.com/boscod/responsewatch-cli/internal/models"
)

// Client represents the API client
type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Config     *config.Config
}

// NewClient creates a new API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Duration(cfg.API.Timeout) * time.Second,
		},
		BaseURL: cfg.API.BaseURL,
		Config:  cfg,
	}
}

// SetBaseURL sets a custom base URL
func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

// doRequest makes an HTTP request
func (c *Client) doRequest(method, path string, body interface{}, auth bool) (*http.Response, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "rwcli/1.0")

	if auth && c.Config.Auth.Token != "" {
		// Check if token needs refresh
		if c.Config.Auth.ExpiresAt.Before(time.Now().Add(5 * time.Minute)) && c.Config.Auth.RefreshToken != "" {
			if err := c.refreshToken(); err != nil {
				return nil, fmt.Errorf("token expired and refresh failed: %w", err)
			}
		}
		req.Header.Set("Authorization", "Bearer "+c.Config.Auth.Token)
	}

	return c.HTTPClient.Do(req)
}

// refreshToken refreshes the access token
func (c *Client) refreshToken() error {
	url := c.BaseURL + "/api/auth/refresh"

	reqBody := models.RefreshTokenRequest{
		RefreshToken: c.Config.Auth.RefreshToken,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh token failed with status %d", resp.StatusCode)
	}

	var refreshResp models.RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return err
	}

	// Update config
	c.Config.Auth.Token = refreshResp.AccessToken
	c.Config.Auth.RefreshToken = refreshResp.RefreshToken
	c.Config.Auth.ExpiresAt = time.Now().Add(time.Duration(refreshResp.ExpiresIn) * time.Second)

	return c.Config.Save()
}

// decodeResponse decodes the response body
func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil {
			return fmt.Errorf("%s (HTTP %d)", errResp.Message, resp.StatusCode)
		}
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	if v != nil && len(body) > 0 {
		if err := json.Unmarshal(body, v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Get makes a GET request
func (c *Client) Get(path string, v interface{}, auth bool) error {
	resp, err := c.doRequest(http.MethodGet, path, nil, auth)
	if err != nil {
		return err
	}
	return decodeResponse(resp, v)
}

// Post makes a POST request
func (c *Client) Post(path string, body, v interface{}, auth bool) error {
	resp, err := c.doRequest(http.MethodPost, path, body, auth)
	if err != nil {
		return err
	}
	return decodeResponse(resp, v)
}

// Put makes a PUT request
func (c *Client) Put(path string, body, v interface{}, auth bool) error {
	resp, err := c.doRequest(http.MethodPut, path, body, auth)
	if err != nil {
		return err
	}
	return decodeResponse(resp, v)
}

// Patch makes a PATCH request
func (c *Client) Patch(path string, body, v interface{}, auth bool) error {
	resp, err := c.doRequest(http.MethodPatch, path, body, auth)
	if err != nil {
		return err
	}
	return decodeResponse(resp, v)
}

// Delete makes a DELETE request
func (c *Client) Delete(path string, v interface{}, auth bool) error {
	resp, err := c.doRequest(http.MethodDelete, path, nil, auth)
	if err != nil {
		return err
	}
	return decodeResponse(resp, v)
}
