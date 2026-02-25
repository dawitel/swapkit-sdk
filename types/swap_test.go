package types

import (
	"encoding/json"
	"testing"
)

func TestSwapRequest_JSON(t *testing.T) {
	req := &SwapRequest{
		RouteID:            "r1",
		SourceAddress:      "0xsrc",
		DestinationAddress: "0xdst",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var decoded SwapRequest
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.RouteID != req.RouteID || decoded.SourceAddress != req.SourceAddress {
		t.Errorf("decoded: %+v", decoded)
	}
}

func TestSwapResponse_TxAsString(t *testing.T) {
	// API can return tx as base64 string (e.g. PSBT).
	j := `{"routeId":"r1","targetAddress":"0xabc","inboundAddress":"0xdef","tx":"cHNidP8BAHUCAAAAAQ==","meta":{"txType":"PSBT"}}`
	var resp SwapResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.RouteID != "r1" || resp.TargetAddress != "0xabc" {
		t.Errorf("resp: %+v", resp)
	}
	if len(resp.Tx) == 0 {
		t.Error("tx should be decoded")
	}
	if resp.Meta == nil || resp.Meta.TxType != "PSBT" {
		t.Errorf("meta: %+v", resp.Meta)
	}
}

func TestSwapResponse_TxAsObject(t *testing.T) {
	// API can return tx as object (e.g. EVM transaction).
	j := `{"routeId":"r1","tx":{"to":"0x4e69","from":"0xd8dA","gas":"0x5208","value":"0x29a2241d62dea8","data":"0x"},"meta":{"txType":"EVM"}}`
	var resp SwapResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Tx) == 0 {
		t.Error("tx object should be stored as RawMessage")
	}
	if resp.Meta == nil || resp.Meta.TxType != "EVM" {
		t.Errorf("meta: %+v", resp.Meta)
	}
}
