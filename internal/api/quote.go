package api

import (
	"context"

	"github.com/dawitel/swapkit-sdk/types"
)

// GetQuote requests a swap quote (POST /v3/quote).
func GetQuote(ctx context.Context, c Doer, req *types.QuoteRequest) (*types.QuoteResponse, error) {
	var out types.QuoteResponse
	if err := c.Do(ctx, "POST", "/v3/quote", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
