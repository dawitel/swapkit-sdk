package types

import (
	"encoding/json"
	"testing"
)

func TestPriceRequest_JSON(t *testing.T) {
	req := PriceRequest{
		Tokens:   []PriceTokenInput{{Identifier: "ETH.ETH"}, {Identifier: "BTC.BTC"}},
		Metadata: true,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var decoded PriceRequest
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.Tokens) != 2 || decoded.Tokens[0].Identifier != "ETH.ETH" || !decoded.Metadata {
		t.Errorf("decoded: %+v", decoded)
	}
}

func TestPriceResult_WithCoingecko(t *testing.T) {
	j := `{"identifier":"ETH.ETH","provider":"","price_usd":1653.61,"timestamp":1744720254562,"cg":{"id":"ethereum","name":"Ethereum","market_cap":197138665861,"total_volume":12560864823,"price_change_24h_usd":-39.89,"price_change_percentage_24h_usd":-2.38,"timestamp":"2025-04-15T12:30:44.643Z"}}`
	var r PriceResult
	if err := json.Unmarshal([]byte(j), &r); err != nil {
		t.Fatal(err)
	}
	if r.Identifier != "ETH.ETH" || r.PriceUSD != 1653.61 || r.Timestamp != 1744720254562 {
		t.Errorf("PriceResult: %+v", r)
	}
	if r.Cg == nil || r.Cg.ID != "ethereum" || r.Cg.Name != "Ethereum" {
		t.Errorf("Cg: %+v", r.Cg)
	}
	if r.Cg.PriceChange24hUSD != -39.89 || r.Cg.PriceChangePct24hUSD != -2.38 {
		t.Errorf("Cg fields: %+v", r.Cg)
	}
}

func TestPriceResponse_Array(t *testing.T) {
	j := `[{"identifier":"ETH.ETH","price_usd":2000,"timestamp":1},{"identifier":"BTC.BTC","price_usd":50000,"timestamp":2}]`
	var resp PriceResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Fatalf("len(resp)=%d", len(resp))
	}
	if resp[0].Identifier != "ETH.ETH" || resp[0].PriceUSD != 2000 {
		t.Errorf("resp[0]: %+v", resp[0])
	}
	if resp[1].Identifier != "BTC.BTC" || resp[1].PriceUSD != 50000 {
		t.Errorf("resp[1]: %+v", resp[1])
	}
}
