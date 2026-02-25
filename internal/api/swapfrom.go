package api

import (
	"context"
	"net/url"
)

// GetSwapFrom returns token identifiers that can be sold to receive the given buy asset (GET /swapFrom?buyAsset=...).
func GetSwapFrom(ctx context.Context, c Doer, buyAsset string) ([]string, error) {
	q := url.Values{}
	q.Set("buyAsset", buyAsset)
	var out []string
	if err := c.Do(ctx, "GET", "/swapFrom", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
