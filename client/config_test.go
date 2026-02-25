package client

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig("mykey")
	if cfg.APIKey != "mykey" {
		t.Errorf("APIKey: %s", cfg.APIKey)
	}
	if cfg.BaseURL != DefaultBaseURL {
		t.Errorf("BaseURL: %s", cfg.BaseURL)
	}
	if cfg.Timeout != DefaultTimeout {
		t.Errorf("Timeout: %v", cfg.Timeout)
	}
	if DefaultTimeout != 30*time.Second {
		t.Errorf("DefaultTimeout: %v", DefaultTimeout)
	}
}
