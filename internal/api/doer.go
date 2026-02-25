package api

import (
	"context"
	"net/url"
)

// Doer performs HTTP requests to the SwapKit API. *client.Client implements this.
type Doer interface {
	Do(ctx context.Context, method, path string, query url.Values, body interface{}, result interface{}) error
}
