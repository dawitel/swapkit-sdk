package client

import (
	"context"

	"github.com/dawitel/swapkit-sdk/internal/api"
	"github.com/dawitel/swapkit-sdk/types"
)

// Providers returns supported chains per provider (GET /providers).
func (c *Client) Providers(ctx context.Context) ([]types.Provider, error) {
	return api.GetProviders(ctx, c)
}

// Tokens returns supported tokens for a provider (GET /tokens).
func (c *Client) Tokens(ctx context.Context, provider string) (*types.TokenList, error) {
	return api.GetTokens(ctx, c, provider)
}

// SwapFrom returns token identifiers that can be sold to receive buyAsset (GET /swapFrom).
func (c *Client) SwapFrom(ctx context.Context, buyAsset string) ([]string, error) {
	return api.GetSwapFrom(ctx, c, buyAsset)
}

// SwapTo returns token identifiers that can be bought with sellAsset (GET /swapTo).
func (c *Client) SwapTo(ctx context.Context, sellAsset string) ([]string, error) {
	return api.GetSwapTo(ctx, c, sellAsset)
}

// Quote requests a swap quote (POST /v3/quote).
func (c *Client) Quote(ctx context.Context, req *types.QuoteRequest) (*types.QuoteResponse, error) {
	return api.GetQuote(ctx, c, req)
}

// Swap builds swap transaction for a chosen route (POST /v3/swap).
func (c *Client) Swap(ctx context.Context, req *types.SwapRequest) (*types.SwapResponse, error) {
	return api.GetSwap(ctx, c, req)
}

// Track returns status of a swap (POST /track).
func (c *Client) Track(ctx context.Context, req *types.TrackRequest) (*types.TrackResponse, error) {
	return api.GetTrack(ctx, c, req)
}

// Price returns token prices (POST /price).
func (c *Client) Price(ctx context.Context, req *types.PriceRequest) (types.PriceResponse, error) {
	return api.GetPrice(ctx, c, req)
}
