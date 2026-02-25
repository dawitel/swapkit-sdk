package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dawitel/swapkit-sdk/types"
)

func TestClient_Providers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"name":"THORCHAIN","provider":"THORCHAIN","count":5}]`))
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	list, err := c.Providers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Name != "THORCHAIN" {
		t.Errorf("list: %+v", list)
	}
}

func TestClient_Quote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"quoteId":"q1","routes":[{"routeId":"r1"}]}`))
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	resp, err := c.Quote(context.Background(), &types.QuoteRequest{SellAsset: "ETH.ETH", BuyAsset: "BTC.BTC", SellAmount: "1", Slippage: 2})
	if err != nil {
		t.Fatal(err)
	}
	if resp.QuoteID != "q1" || len(resp.Routes) != 1 || resp.Routes[0].RouteID != "r1" {
		t.Errorf("resp: %+v", resp)
	}
}

func TestClient_Price(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"identifier":"ETH.ETH","price_usd":2000,"timestamp":1}]`))
	}))
	defer srv.Close()
	c, _ := New(Config{BaseURL: srv.URL, APIKey: "k", HTTPClient: srv.Client()})
	resp, err := c.Price(context.Background(), &types.PriceRequest{Tokens: []types.PriceTokenInput{{Identifier: "ETH.ETH"}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) != 1 || resp[0].Identifier != "ETH.ETH" || resp[0].PriceUSD != 2000 {
		t.Errorf("resp: %v", resp)
	}
}
