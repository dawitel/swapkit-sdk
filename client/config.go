package client

import (
	"net/http"
	"time"
)

const DefaultBaseURL = "https://api.swapkit.dev"
const DefaultTimeout = 30 * time.Second

// Config holds client configuration.
type Config struct {
	BaseURL    string
	APIKey     string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// DefaultConfig returns config with default base URL and timeout; APIKey must be set by caller.
func DefaultConfig(apiKey string) Config {
	return Config{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		Timeout: DefaultTimeout,
	}
}
