package types

// PriceRequest is the body for POST /price.
type PriceRequest struct {
	Tokens   []PriceTokenInput `json:"tokens"`
	Metadata bool              `json:"metadata,omitempty"`
}

// PriceTokenInput is one token identifier for price lookup.
type PriceTokenInput struct {
	Identifier string `json:"identifier"`
}

// PriceResponse is the array returned by POST /price.
type PriceResponse []PriceResult

// PriceResult is one token price.
type PriceResult struct {
	Identifier string          `json:"identifier"`
	Provider   string          `json:"provider,omitempty"`
	PriceUSD   float64         `json:"price_usd"`
	Timestamp  int64           `json:"timestamp"`
	Cg         *CoingeckoPrice `json:"cg,omitempty"`
}

// CoingeckoPrice is CoinGecko metadata when metadata=true.
type CoingeckoPrice struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	MarketCap            float64   `json:"market_cap,omitempty"`
	TotalVolume          float64   `json:"total_volume,omitempty"`
	PriceChange24hUSD    float64   `json:"price_change_24h_usd,omitempty"`
	PriceChangePct24hUSD float64   `json:"price_change_percentage_24h_usd,omitempty"`
	SparklineIn7d        []float64 `json:"sparkline_in_7d,omitempty"`
	Timestamp            string    `json:"timestamp,omitempty"`
}
