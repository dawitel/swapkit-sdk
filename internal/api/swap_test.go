package api_test

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dawitel/swapkit-sdk/client"
	swapkitapi "github.com/dawitel/swapkit-sdk/internal/api"
	"github.com/dawitel/swapkit-sdk/types"
)

func TestGetSwap(t *testing.T) {
	txB64 := base64.StdEncoding.EncodeToString([]byte("psbt-data"))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v3/swap" || r.Method != "POST" {
			t.Errorf("path=%s method=%s", r.URL.Path, r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"routeId":"r1","targetAddress":"0xabc","inboundAddress":"0xdef","tx":"` + txB64 + `","meta":{"txType":"PSBT"}}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	req := &types.SwapRequest{RouteID: "r1", SourceAddress: "0xs", DestinationAddress: "0xd"}
	resp, err := swapkitapi.GetSwap(context.Background(), c, req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.RouteID != "r1" || resp.TargetAddress != "0xabc" || resp.InboundAddress != "0xdef" {
		t.Errorf("resp: %+v", resp)
	}
	if len(resp.Tx) == 0 {
		t.Error("tx empty")
	}
	if resp.Meta == nil || resp.Meta.TxType != "PSBT" {
		t.Errorf("meta: %+v", resp.Meta)
	}
}

func TestGetSwap_Error(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"insufficientBalance","message":"Insufficient balance"}`))
	}))
	defer srv.Close()
	c := newTestClient(t, srv)
	_, err := swapkitapi.GetSwap(context.Background(), c, &types.SwapRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, client.ErrInsufficientBalance) {
		t.Errorf("expected ErrInsufficientBalance: %v", err)
	}
}
