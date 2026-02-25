package types

import (
	"encoding/json"
	"testing"
)

func TestProvider_JSONRoundTrip(t *testing.T) {
	p := Provider{
		Name:              "THORCHAIN",
		Provider:          "THORCHAIN",
		Keywords:          []string{"swap"},
		Count:             10,
		LogoURI:           "https://example.com/logo.png",
		URL:               "https://example.com/tokens.json",
		SupportedActions:  []string{"swap"},
		SupportedChainIds: []string{"bitcoin", "1", "solana"},
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Provider
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Name != p.Name || decoded.Provider != p.Provider || decoded.Count != p.Count {
		t.Errorf("decoded: %+v", decoded)
	}
	if len(decoded.SupportedChainIds) != 3 {
		t.Errorf("SupportedChainIds: %v", decoded.SupportedChainIds)
	}
}

func TestProvider_UnmarshalFromAPI(t *testing.T) {
	j := `{"name":"CHAINFLIP_STREAMING","provider":"CHAINFLIP_STREAMING","keywords":[],"count":10,"logoURI":"https://example.com/logo.png","url":"https://example.com/list.json","supportedActions":["swap"],"supportedChainIds":["42161","bitcoin","1","solana"]}`
	var p Provider
	if err := json.Unmarshal([]byte(j), &p); err != nil {
		t.Fatal(err)
	}
	if p.Name != "CHAINFLIP_STREAMING" || p.Count != 10 || len(p.SupportedChainIds) != 4 {
		t.Errorf("Provider: %+v", p)
	}
}
