package api

import (
	"context"
	"net/url"

	"github.com/dawitel/swapkit-sdk/types"
)

// GetTokens returns supported tokens for a provider (GET /tokens?provider=...).
func GetTokens(ctx context.Context, c Doer, provider string) (*types.TokenList, error) {
	q := url.Values{}
	q.Set("provider", provider)
	var out types.TokenList
	if err := c.Do(ctx, "GET", "/tokens", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
