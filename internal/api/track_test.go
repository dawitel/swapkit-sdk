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

func TestGetTrack(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/track" || r.Method != "POST" {
			t.Errorf("path=%s method=%s", r.URL.Path, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"chainId":"1","hash":"0xabc","status":"completed","trackingStatus":"completed","fromAsset":"ETH.ETH","toAsset":"BTC.BTC"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	req := &types.TrackRequest{Hash: "0xabc", ChainID: "1"}
	resp, err := swapkitapi.GetTrack(context.Background(), c, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ChainID != "1" || resp.Hash != "0xabc" || resp.Status != "completed" {
		t.Errorf("resp: %+v", resp)
	}
	if resp.FromAsset != "ETH.ETH" || resp.ToAsset != "BTC.BTC" {
		t.Errorf("assets: %s -> %s", resp.FromAsset, resp.ToAsset)
	}
}

func TestGetTrack_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalidRequest","message":"hash or depositAddress required"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetTrack(context.Background(), c, &types.TrackRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError: %T", err)
	}
}
