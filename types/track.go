package types

// TrackRequest is the body for POST /track (hash+chainId or depositAddress).
type TrackRequest struct {
	Hash           string `json:"hash,omitempty"`
	ChainID        string `json:"chainId,omitempty"`
	DepositAddress string `json:"depositAddress,omitempty"`
}

// TrackResponse is the response from POST /track.
type TrackResponse struct {
	ChainID        string                 `json:"chainId"`
	Hash           string                 `json:"hash"`
	Block          int64                  `json:"block"`
	Type           string                 `json:"type"`
	Status         string                 `json:"status"`
	TrackingStatus string                 `json:"trackingStatus"`
	FromAsset      string                 `json:"fromAsset"`
	FromAmount     string                 `json:"fromAmount"`
	FromAddress    string                 `json:"fromAddress"`
	ToAsset        string                 `json:"toAsset"`
	ToAmount       string                 `json:"toAmount"`
	ToAddress      string                 `json:"toAddress"`
	FinalisedAt    int64                  `json:"finalisedAt"`
	Meta           *TrackMeta             `json:"meta,omitempty"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
	Legs           []TrackResponse        `json:"legs,omitempty"`
}

// TrackMeta holds provider and image URLs.
type TrackMeta struct {
	Provider       string            `json:"provider,omitempty"`
	ProviderAction string            `json:"providerAction,omitempty"`
	Images         map[string]string `json:"images,omitempty"`
}
