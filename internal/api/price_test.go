package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dawitel/swapkit-sdk/client"
	swapkitapi "github.com/dawitel/swapkit-sdk/internal/api"
	"github.com/dawitel/swapkit-sdk/types"
)

func TestGetPrice(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/price" || r.Method != "POST" {
			t.Errorf("path=%s method=%s", r.URL.Path, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"identifier":"ETH.ETH","price_usd":2500,"timestamp":1234567890}]`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	req := &types.PriceRequest{Tokens: []types.PriceTokenInput{{Identifier: "ETH.ETH"}}}
	resp, err := swapkitapi.GetPrice(context.Background(), c, req)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) != 1 {
		t.Fatalf("len(resp)=%d", len(resp))
	}
	if resp[0].Identifier != "ETH.ETH" || resp[0].PriceUSD != 2500 || resp[0].Timestamp != 1234567890 {
		t.Errorf("result: %+v", resp[0])
	}
}

func TestGetPrice_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"apiKeyInvalid"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetPrice(context.Background(), c, &types.PriceRequest{Tokens: []types.PriceTokenInput{{Identifier: "ETH.ETH"}}})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError: %T", err)
	}
}
