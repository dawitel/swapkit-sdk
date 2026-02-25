package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dawitel/swapkit-sdk/client"
	swapkitapi "github.com/dawitel/swapkit-sdk/internal/api"
)

func TestGetTokens(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tokens" || r.URL.Query().Get("provider") != "THORCHAIN" {
			t.Errorf("path=%s query=%s", r.URL.Path, r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"provider":"THORCHAIN","count":2,"tokens":[{"chain":"BTC","chainId":"bitcoin","identifier":"BTC.BTC","ticker":"BTC"}]}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	list, err := swapkitapi.GetTokens(context.Background(), c, "THORCHAIN")
	if err != nil {
		t.Fatal(err)
	}
	if list.Provider != "THORCHAIN" || list.Count != 2 || len(list.Tokens) != 1 {
		t.Errorf("list: %+v", list)
	}
	if list.Tokens[0].Identifier != "BTC.BTC" {
		t.Errorf("token: %+v", list.Tokens[0])
	}
}

func TestGetTokens_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalidRequest","message":"provider required"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetTokens(context.Background(), c, "X")
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError: %T", err)
	}
}
