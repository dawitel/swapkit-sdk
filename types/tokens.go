package types

// TokenList is the response from GET /tokens?provider=...
type TokenList struct {
	Provider  string   `json:"provider"`
	Name      string   `json:"name"`
	Timestamp string   `json:"timestamp,omitempty"`
	Version   Version  `json:"version,omitempty"`
	Keywords  []string `json:"keywords,omitempty"`
	Count     int      `json:"count"`
	Tokens    []Token  `json:"tokens"`
}

// Version is token list version.
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

// Token is a supported token (identifier is used in quote/swap).
type Token struct {
	Chain       string `json:"chain"`
	ChainID     string `json:"chainId"`
	Address     string `json:"address,omitempty"`
	Ticker      string `json:"ticker"`
	Identifier  string `json:"identifier"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Decimals    int    `json:"decimals"`
	LogoURI     string `json:"logoURI,omitempty"`
	CoingeckoID string `json:"coingeckoId,omitempty"`
}
