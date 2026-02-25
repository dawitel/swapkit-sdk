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

func TestGetProviders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/providers" || r.Method != "GET" {
			t.Errorf("path=%s method=%s", r.URL.Path, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"THORCHAIN","provider":"THORCHAIN","supportedChainIds":["bitcoin","1"],"count":10}]`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	list, err := swapkitapi.GetProviders(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len(list)=%d", len(list))
	}
	if list[0].Name != "THORCHAIN" || list[0].Provider != "THORCHAIN" || list[0].Count != 10 {
		t.Errorf("item: %+v", list[0])
	}
	if len(list[0].SupportedChainIds) != 2 {
		t.Errorf("supportedChainIds: %v", list[0].SupportedChainIds)
	}
}

func TestGetProviders_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"apiKeyInvalid","message":"Invalid API key"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetProviders(context.Background(), c)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError: %T", err)
	}
}
