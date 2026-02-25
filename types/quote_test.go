package types

import (
	"encoding/json"
	"testing"
)

func TestQuoteRequestResponse_JSON(t *testing.T) {
	req := &QuoteRequest{
		SellAsset: "ETH.ETH", BuyAsset: "BTC.BTC", SellAmount: "1", Slippage: 2,
		Providers: []string{"THORCHAIN"},
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var decoded QuoteRequest
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.SellAsset != req.SellAsset || decoded.BuyAsset != req.BuyAsset || len(decoded.Providers) != 1 {
		t.Errorf("decoded: %+v", decoded)
	}

	respJSON := `{"quoteId":"q1","routes":[{"routeId":"r1","outboundAmount":"0.005","estimatedTime":{"total":660},"meta":{"tags":["FASTEST"]}}]}`
	var resp QuoteResponse
	if err := json.Unmarshal([]byte(respJSON), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.QuoteID != "q1" || len(resp.Routes) != 1 {
		t.Errorf("resp: %+v", resp)
	}
	if resp.Routes[0].RouteID != "r1" || resp.Routes[0].OutboundAmount != "0.005" {
		t.Errorf("route: %+v", resp.Routes[0])
	}
	if resp.Routes[0].EstimatedTime == nil || resp.Routes[0].EstimatedTime.Total != 660 {
		t.Errorf("estimatedTime: %+v", resp.Routes[0].EstimatedTime)
	}
	if resp.Routes[0].Meta == nil || len(resp.Routes[0].Meta.Tags) != 1 || resp.Routes[0].Meta.Tags[0] != "FASTEST" {
		t.Errorf("meta: %+v", resp.Routes[0].Meta)
	}
}

func TestQuoteRoute_FeesAsArray(t *testing.T) {
	// API returns fees as an array of objects, not a single object
	j := `{"quoteId":"q1","routes":[{"routeId":"r1","outboundAmount":"0.1","fees":[{"type":"inbound","amount":"0.0001","asset":"ETH.ETH"},{"type":"swap","amount":"0.001","asset":"ETH.ETH"}]}]}`
	var resp QuoteResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Routes) != 1 || len(resp.Routes[0].Fees) == 0 {
		t.Errorf("routes or fees: %+v", resp.Routes[0])
	}
}

func TestRouteByTag(t *testing.T) {
	// Nil or empty response
	if got := RouteByTag(nil, TagCheapest); got != nil {
		t.Errorf("RouteByTag(nil): want nil, got %+v", got)
	}
	if got := RouteByTag(&QuoteResponse{}, TagCheapest); got != nil {
		t.Errorf("RouteByTag(no routes): want nil, got %+v", got)
	}

	// Two routes: first FASTEST, second CHEAPEST
	resp := &QuoteResponse{
		QuoteID: "q1",
		Routes: []QuoteRoute{
			{RouteID: "r1", Meta: &QuoteRouteMeta{Tags: []string{"FASTEST"}}},
			{RouteID: "r2", Meta: &QuoteRouteMeta{Tags: []string{"CHEAPEST"}}},
		},
	}
	if got := RouteByTag(resp, TagCheapest); got == nil || got.RouteID != "r2" {
		t.Errorf("RouteByTag(CHEAPEST): want r2, got %v", got)
	}
	if got := RouteByTag(resp, TagFastest); got == nil || got.RouteID != "r1" {
		t.Errorf("RouteByTag(FASTEST): want r1, got %v", got)
	}
	// No route has RECOMMENDED -> returns first route
	if got := RouteByTag(resp, TagRecommended); got == nil || got.RouteID != "r1" {
		t.Errorf("RouteByTag(RECOMMENDED) fallback: want r1, got %v", got)
	}
}
