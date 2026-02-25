package types

import "encoding/json"

// QuoteRequest is the body for POST /v3/quote.
type QuoteRequest struct {
	SellAsset  string   `json:"sellAsset"`
	BuyAsset   string   `json:"buyAsset"`
	SellAmount string   `json:"sellAmount"`
	Slippage   int      `json:"slippage"`
	Providers  []string `json:"providers,omitempty"`
}

// QuoteResponse is the response from POST /v3/quote.
type QuoteResponse struct {
	QuoteID        string       `json:"quoteId"`
	Routes         []QuoteRoute `json:"routes"`
	ProviderErrors []QuoteError `json:"providerErrors,omitempty"`
	Error          string       `json:"error,omitempty"`
}

// QuoteRoute is one route in a quote.
type QuoteRoute struct {
	RouteID                    string          `json:"routeId"`
	Provider                   string          `json:"provider,omitempty"`
	Providers                  []string        `json:"providers,omitempty"`
	SellAsset                  string          `json:"sellAsset,omitempty"`
	SellAmount                 string          `json:"sellAmount,omitempty"`
	BuyAsset                   string          `json:"buyAsset,omitempty"`
	ExpectedBuyAmount          string          `json:"expectedBuyAmount,omitempty"`
	ExpectedBuyAmountMaxSlippage string        `json:"expectedBuyAmountMaxSlippage,omitempty"`
	InboundAmount              string          `json:"inboundAmount,omitempty"`
	OutboundAmount             string          `json:"outboundAmount,omitempty"`
	Expiration                 string          `json:"expiration,omitempty"`
	EstimatedTime              *EstimatedTime  `json:"estimatedTime,omitempty"`
	Fees                       json.RawMessage `json:"fees,omitempty"` // API returns an array of fee objects
	TotalSlippage              string          `json:"totalSlippage,omitempty"`
	TotalSlippageBps           float64         `json:"totalSlippageBps,omitempty"`
	Legs                       json.RawMessage `json:"legs,omitempty"`
	Warnings                   json.RawMessage `json:"warnings,omitempty"`
	Meta                       *QuoteRouteMeta `json:"meta,omitempty"`
	NextActions                []NextAction    `json:"nextActions,omitempty"`
}

// EstimatedTime is estimated timing for a swap (seconds).
type EstimatedTime struct {
	Inbound  float64 `json:"inbound"`
	Swap     float64 `json:"swap"`
	Outbound float64 `json:"outbound"`
	Total    float64 `json:"total"`
}

// QuoteRouteMeta is metadata for a quote route.
type QuoteRouteMeta struct {
	Assets []AssetMeta `json:"assets,omitempty"`
	Tags   []string    `json:"tags,omitempty"`
}

// AssetMeta is asset info in quote meta.
type AssetMeta struct {
	Asset string  `json:"asset"`
	Price float64 `json:"price"`
	Image string  `json:"image,omitempty"`
}

// NextAction describes the next API call (e.g. /swap).
type NextAction struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Payload interface{} `json:"payload,omitempty"`
}

// QuoteError is a provider-level quote error.
type QuoteError struct {
	Provider  string `json:"provider"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

// Well-known route tags returned by the SwapKit API (meta.tags).
// The API aggregates quotes from multiple providers and marks the best route(s).
const (
	TagRecommended = "RECOMMENDED" // API's default pick (often best overall)
	TagFastest     = "FASTEST"     // Shortest estimated time
	TagCheapest    = "CHEAPEST"    // Best expected buy amount / lowest fees
)

// RouteByTag returns the first route that has the given tag (e.g. TagCheapest, TagFastest, TagRecommended).
// SwapKit returns multiple routes and sets meta.tags so you can pick by preference.
// If no route has the tag, returns the first route (typically the API's recommended one). Returns nil if no routes.
func RouteByTag(resp *QuoteResponse, tag string) *QuoteRoute {
	if resp == nil || len(resp.Routes) == 0 {
		return nil
	}
	for i := range resp.Routes {
		if resp.Routes[i].Meta != nil {
			for _, t := range resp.Routes[i].Meta.Tags {
				if t == tag {
					return &resp.Routes[i]
				}
			}
		}
	}
	return &resp.Routes[0]
}
