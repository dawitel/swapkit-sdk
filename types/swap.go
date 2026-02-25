package types

import "encoding/json"

// SwapRequest is the body for POST /v3/swap.
type SwapRequest struct {
	RouteID             string `json:"routeId"`
	SourceAddress       string `json:"sourceAddress"`
	DestinationAddress  string `json:"destinationAddress"`
	DisableBalanceCheck *bool  `json:"disableBalanceCheck,omitempty"`
	DisableBuildTx      *bool  `json:"disableBuildTx,omitempty"`
	OverrideSlippage    *bool  `json:"overrideSlippage,omitempty"`
	DisableEstimate     *bool  `json:"disableEstimate,omitempty"`
}

// SwapResponse is the response from POST /v3/swap (single route with tx fields).
type SwapResponse struct {
	QuoteRoute
	TargetAddress  string          `json:"targetAddress,omitempty"`
	InboundAddress string          `json:"inboundAddress,omitempty"`
	Tx             json.RawMessage `json:"tx,omitempty"`
	Meta           *SwapMeta       `json:"meta,omitempty"`
}

// SwapMeta can include txType (e.g. "PSBT", "EVM").
type SwapMeta struct {
	QuoteRouteMeta
	TxType string `json:"txType,omitempty"`
}
