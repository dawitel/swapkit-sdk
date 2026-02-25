package api

import (
	"context"
	"net/url"
)

// GetSwapTo returns token identifiers that can be bought with the given sell asset (GET /swapTo?sellAsset=...).
func GetSwapTo(ctx context.Context, c Doer, sellAsset string) ([]string, error) {
	q := url.Values{}
	q.Set("sellAsset", sellAsset)
	var out []string
	if err := c.Do(ctx, "GET", "/swapTo", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
