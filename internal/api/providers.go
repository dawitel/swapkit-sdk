package api

import (
	"context"

	"github.com/dawitel/swapkit-sdk/types"
)

// GetProviders returns supported chains per provider (GET /providers).
func GetProviders(ctx context.Context, c Doer) ([]types.Provider, error) {
	var out []types.Provider
	if err := c.Do(ctx, "GET", "/providers", nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
