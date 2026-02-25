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

func TestGetSwapFrom(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/swapFrom" || r.URL.Query().Get("buyAsset") != "BTC.BTC" {
			t.Errorf("path=%s query=%s", r.URL.Path, r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`["ETH.ETH","AVAX.AVAX"]`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	list, err := swapkitapi.GetSwapFrom(context.Background(), c, "BTC.BTC")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0] != "ETH.ETH" || list[1] != "AVAX.AVAX" {
		t.Errorf("list: %v", list)
	}
}

func TestGetSwapFrom_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"noRoutesFound"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetSwapFrom(context.Background(), c, "INVALID")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, client.ErrNoRoutesFound) {
		t.Errorf("expected ErrNoRoutesFound: %v", err)
	}
}
