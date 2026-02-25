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

func TestGetQuote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/quote" || r.Method != "POST" {
			t.Errorf("path=%s method=%s", r.URL.Path, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"quoteId":"q1","routes":[{"routeId":"r1","outboundAmount":"0.005","meta":{"tags":["RECOMMENDED"]}}]}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	req := &types.QuoteRequest{SellAsset: "ETH.ETH", BuyAsset: "BTC.BTC", SellAmount: "1", Slippage: 2}
	resp, err := swapkitapi.GetQuote(context.Background(), c, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.QuoteID != "q1" || len(resp.Routes) != 1 {
		t.Errorf("resp: %+v", resp)
	}
	if resp.Routes[0].RouteID != "r1" || resp.Routes[0].OutboundAmount != "0.005" {
		t.Errorf("route: %+v", resp.Routes[0])
	}
	if len(resp.Routes[0].Meta.Tags) != 1 || resp.Routes[0].Meta.Tags[0] != "RECOMMENDED" {
		t.Errorf("meta.tags: %v", resp.Routes[0].Meta.Tags)
	}
}

func TestGetQuote_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"noRoutesFound","message":"No routes found for X -> Y"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetQuote(context.Background(), c, &types.QuoteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError: %T", err)
	}
	if !errors.Is(err, client.ErrNoRoutesFound) {
		t.Errorf("expected ErrNoRoutesFound: %v", err)
	}
}
