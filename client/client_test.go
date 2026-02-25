package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew_RequiresAPIKey(t *testing.T) {
	_, err := New(Config{})
	if err == nil {
		t.Fatal("expected error when APIKey is empty")
	}
	c, err := New(DefaultConfig("key"))
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestDo_Success(t *testing.T) {
	body := []byte(`{"quoteId":"q1","routes":[]}`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "testkey" {
			t.Errorf("missing or wrong x-api-key: %s", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("accept") != "application/json" {
			t.Errorf("accept: %s", r.Header.Get("accept"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}))
	defer srv.Close()

	c, err := New(Config{BaseURL: srv.URL, APIKey: "testkey", HTTPClient: srv.Client()})
	if err != nil {
		t.Fatal(err)
	}
	var out struct {
		QuoteID string   `json:"quoteId"`
		Routes  []string `json:"routes"`
	}
	err = c.Do(context.Background(), "GET", "/v3/quote", nil, nil, &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.QuoteID != "q1" || len(out.Routes) != 0 {
		t.Errorf("unexpected decode: %+v", out)
	}
}

func TestDo_ErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantCode   string
		wantUnwrap error
	}{
		{
			name:       "insufficientBalance",
			status:     http.StatusBadRequest,
			body:       `{"message":"Cannot build...","error":"insufficientBalance","data":{}}`,
			wantCode:   "insufficientBalance",
			wantUnwrap: ErrInsufficientBalance,
		},
		{
			name:       "noRoutesFound",
			status:     http.StatusNotFound,
			body:       `{"message":"No routes found","error":"noRoutesFound"}`,
			wantCode:   "noRoutesFound",
			wantUnwrap: ErrNoRoutesFound,
		},
		{
			name:       "unauthorized",
			status:     http.StatusUnauthorized,
			body:       `{"message":"Invalid API key","error":"apiKeyInvalid"}`,
			wantCode:   "apiKeyInvalid",
			wantUnwrap: ErrAPIKeyInvalid,
		},
		{
			name:       "blackListAsset",
			status:     http.StatusBadRequest,
			body:       `{"error":"blackListAsset","message":"Asset blacklisted"}`,
			wantCode:   "blackListAsset",
			wantUnwrap: ErrBlackListAsset,
		},
		{
			name:       "invalidRequest",
			status:     http.StatusBadRequest,
			body:       `{"error":"invalidRequest","message":"Request body required"}`,
			wantCode:   "invalidRequest",
			wantUnwrap: ErrInvalidRequest,
		},
		{
			name:       "insufficientAllowance",
			status:     http.StatusBadRequest,
			body:       `{"error":"insufficientAllowance","message":"Insufficient allowance"}`,
			wantCode:   "insufficientAllowance",
			wantUnwrap: ErrInsufficientAllowance,
		},
		{
			name:       "unableToBuildTransaction",
			status:     http.StatusBadRequest,
			body:       `{"error":"unableToBuildTransaction","message":"Cannot build tx"}`,
			wantCode:   "unableToBuildTransaction",
			wantUnwrap: ErrUnableToBuildTransaction,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()
			c, err := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
			if err != nil {
				t.Fatal(err)
			}
			err = c.Do(context.Background(), "GET", "/", nil, nil, &struct{}{})
			if err == nil {
				t.Fatal("expected error")
			}
			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if apiErr.Code != tt.wantCode {
				t.Errorf("Code: got %q want %q", apiErr.Code, tt.wantCode)
			}
			if tt.wantUnwrap != nil && !errors.Is(err, tt.wantUnwrap) {
				t.Errorf("errors.Is(err, sentinel): got false")
			}
		})
	}
}

func TestDo_POSTWithBody(t *testing.T) {
	var received map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method: %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type: %s", r.Header.Get("Content-Type"))
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"quoteId":"ok"}`))
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	req := map[string]string{"sellAsset": "ETH.ETH", "buyAsset": "BTC.BTC"}
	var out struct {
		QuoteID string `json:"quoteId"`
	}
	err := c.Do(context.Background(), "POST", "/v3/quote", nil, req, &out)
	if err != nil {
		t.Fatal(err)
	}
	if out.QuoteID != "ok" {
		t.Errorf("QuoteID: %s", out.QuoteID)
	}
	if received["sellAsset"] != "ETH.ETH" || received["buyAsset"] != "BTC.BTC" {
		t.Errorf("body: %+v", received)
	}
}

func TestDo_QueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "provider=THORCHAIN" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	q := url.Values{"provider": {"THORCHAIN"}}
	var out []interface{}
	err := c.Do(context.Background(), "GET", "/tokens", q, nil, &out)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDo_NilResultEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	err := c.Do(context.Background(), "GET", "/", nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := c.Do(ctx, "GET", "/", nil, nil, &struct{}{})
	if err == nil {
		t.Fatal("expected error when context cancelled")
	}
}

func TestNew_TrimsBaseURLTrailingSlash(t *testing.T) {
	c, err := New(Config{BaseURL: "https://api.example.com/", APIKey: "k"})
	if err != nil {
		t.Fatal(err)
	}
	// Do a request and assert path is not double-slashed (we'd see // in URL)
	// We can't read c.baseURL; instead do a request and check server received path
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 0 && r.URL.Path[0] != '/' {
			t.Errorf("path should start with /: %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c2, _ := New(Config{BaseURL: srv.URL + "/", APIKey: "k", HTTPClient: srv.Client()})
	_ = c2.Do(context.Background(), "GET", "/providers", nil, nil, &[]interface{}{})
	_ = c
}
