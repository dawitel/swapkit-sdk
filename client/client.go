package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Doer performs HTTP requests. *http.Client implements this.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the SwapKit API HTTP client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient Doer
}

// New creates a Client from config. APIKey is required.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, errors.New("client: APIKey is required")
	}
	base := strings.TrimSuffix(cfg.BaseURL, "/")
	hc := cfg.HTTPClient
	if hc == nil {
		hc = &http.Client{Timeout: cfg.Timeout}
	}
	return &Client{baseURL: base, apiKey: cfg.APIKey, httpClient: hc}, nil
}

// Do sends an HTTP request and decodes the JSON response into result (if not nil).
// On 4xx/5xx, Do returns an *APIError; result is unchanged.
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	slurp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return parseAPIError(resp, slurp)
	}
	if result != nil && len(slurp) > 0 {
		return json.Unmarshal(slurp, result)
	}
	return nil
}
