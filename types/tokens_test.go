package types

import (
	"encoding/json"
	"testing"
)

func TestTokenList_JSONRoundTrip(t *testing.T) {
	list := TokenList{
		Provider:  "THORCHAIN",
		Name:      "THORCHAIN",
		Timestamp: "2025-01-11T16:31:04.355Z",
		Version:   Version{Major: 1, Minor: 0, Patch: 0},
		Count:     2,
		Tokens: []Token{
			{Chain: "BTC", ChainID: "bitcoin", Ticker: "BTC", Identifier: "BTC.BTC", Symbol: "BTC", Name: "Bitcoin", Decimals: 8},
			{Chain: "ETH", ChainID: "1", Ticker: "ETH", Identifier: "ETH.ETH", Symbol: "ETH", Name: "Ether", Decimals: 18},
		},
	}
	b, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	var decoded TokenList
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Provider != list.Provider || decoded.Count != 2 || len(decoded.Tokens) != 2 {
		t.Errorf("decoded: %+v", decoded)
	}
	if decoded.Tokens[0].Identifier != "BTC.BTC" || decoded.Tokens[1].Identifier != "ETH.ETH" {
		t.Errorf("tokens: %+v", decoded.Tokens)
	}
}

func TestToken_WithAddress(t *testing.T) {
	j := `{"chain":"SOL","address":"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","chainId":"solana","ticker":"USDC","identifier":"SOL.USDC-EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","symbol":"USDC-EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v","name":"Solana USDC","decimals":6}`
	var tok Token
	if err := json.Unmarshal([]byte(j), &tok); err != nil {
		t.Fatal(err)
	}
	if tok.Chain != "SOL" || tok.Address == "" || tok.Identifier == "" || tok.Decimals != 6 {
		t.Errorf("Token: %+v", tok)
	}
}
