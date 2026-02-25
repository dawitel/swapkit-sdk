package api

import (
	"context"

	"github.com/dawitel/swapkit-sdk/types"
)

// GetSwap builds swap transaction for a chosen route (POST /v3/swap).
func GetSwap(ctx context.Context, c Doer, req *types.SwapRequest) (*types.SwapResponse, error) {
	var out types.SwapResponse
	if err := c.Do(ctx, "POST", "/v3/swap", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
