package api

import (
	"context"

	"github.com/dawitel/swapkit-sdk/types"
)

// GetTrack returns status of a swap (POST /track).
func GetTrack(ctx context.Context, c Doer, req *types.TrackRequest) (*types.TrackResponse, error) {
	var out types.TrackResponse
	if err := c.Do(ctx, "POST", "/track", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
