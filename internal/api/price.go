package api

import (
	"context"

	"github.com/dawitel/swapkit-sdk/types"
)

// GetPrice returns token prices (POST /price).
func GetPrice(ctx context.Context, c Doer, req *types.PriceRequest) (types.PriceResponse, error) {
	var out types.PriceResponse
	if err := c.Do(ctx, "POST", "/price", nil, req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
