package types

import (
	"encoding/json"
	"testing"
)

func TestTrackRequestResponse_JSON(t *testing.T) {
	req := &TrackRequest{Hash: "0xabc", ChainID: "1"}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var decoded TrackRequest
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Hash != req.Hash || decoded.ChainID != req.ChainID {
		t.Errorf("decoded: %+v", decoded)
	}

	j := `{"chainId":"43114","hash":"0x18f6","block":57181100,"type":"swap","status":"completed","trackingStatus":"completed","fromAsset":"AVAX.AVAX","fromAmount":"9.58","toAsset":"THOR.RUNE","toAmount":"0","meta":{"provider":"THORCHAIN"},"legs":[]}`
	var resp TrackResponse
	if err := json.Unmarshal([]byte(j), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ChainID != "43114" || resp.Status != "completed" || resp.FromAsset != "AVAX.AVAX" || resp.ToAsset != "THOR.RUNE" {
		t.Errorf("resp: %+v", resp)
	}
	if resp.Meta == nil || resp.Meta.Provider != "THORCHAIN" {
		t.Errorf("meta: %+v", resp.Meta)
	}
}
