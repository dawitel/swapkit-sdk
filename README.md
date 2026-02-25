# SwapKit Go SDK

Production-ready Go client for the [SwapKit API](https://docs.swapkit.dev/), covering providers, tokens, swap discovery, v3 quote/swap, track, and price endpoints.

## Install

```bash
go get github.com/dawitel/swapkit-sdk
```

## Configuration

Obtain an API key from the [SwapKit dashboard](https://dashboard.swapkit.dev/). Do not hardcode it; use environment variables or a secrets manager.

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/dawitel/swapkit-sdk/client"
    "github.com/dawitel/swapkit-sdk/types"
)

func main() {
    apiKey := os.Getenv("SWAPKIT_API_KEY")
    if apiKey == "" {
        log.Fatal("SWAPKIT_API_KEY is required")
    }
    cfg := client.DefaultConfig(apiKey)
    // Optional: cfg.BaseURL = "https://dev-api.swapkit.dev"
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    ctx := context.Background()

    // List providers and supported chains
    providers, err := c.Providers(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Get supported tokens for a provider
    tokens, err := c.Tokens(ctx, "THORCHAIN")
    if err != nil {
        log.Fatal(err)
    }

    // Quote (price discovery) – no addresses needed
    quote, err := c.Quote(ctx, &types.QuoteRequest{
        SellAsset:  "ETH.ETH",
        BuyAsset:   "BTC.BTC",
        SellAmount: "1",
        Slippage:   2,
    })
    if err != nil {
        log.Fatal(err)
    }
    if len(quote.Routes) == 0 {
        log.Fatal("no routes")
    }
    routeID := quote.Routes[0].RouteID

    // Build swap transaction (after user accepts the quote)
    swap, err := c.Swap(ctx, &types.SwapRequest{
        RouteID:             routeID,
        SourceAddress:       "0x...",   // user wallet
        DestinationAddress:  "bc1q...", // destination for bought asset
    })
    if err != nil {
        log.Fatal(err)
    }
    // Sign and broadcast swap.Tx (format depends on meta.TxType: EVM, PSBT, etc.)

    // Track status
    track, err := c.Track(ctx, &types.TrackRequest{Hash: "0x...", ChainID: "1"})
    if err != nil {
        log.Fatal(err)
    }

    // Token prices
    prices, err := c.Price(ctx, &types.PriceRequest{
        Tokens:   []types.PriceTokenInput{{Identifier: "ETH.ETH"}, {Identifier: "BTC.BTC"}},
        Metadata: true,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Quote and swap flow

1. **Quote** – `POST /v3/quote` with `sellAsset`, `buyAsset`, `sellAmount`, `slippage`. No addresses. Returns `routeId` per route (cached ~60s).
2. **Swap** – After the user picks a route, call `POST /v3/swap` with the chosen `routeId`, `sourceAddress`, and `destinationAddress`. Response includes `tx` (ready to sign) and `meta.TxType` (EVM, PSBT, etc.).

See [Quote and swap implementation flow](https://docs.swapkit.dev/swapkit-api/quote-and-swap-implementation-flow) in the SwapKit docs.

## Error handling

The client returns `*client.APIError` on 4xx/5xx. You can use `errors.Is` with sentinel errors:

```go
if err != nil {
    if errors.Is(err, client.ErrInsufficientBalance) {
        // handle insufficient balance
    }
    if errors.Is(err, client.ErrNoRoutesFound) {
        // no route for this pair
    }
    if apiErr := (*client.APIError)(nil); errors.As(err, &apiErr) {
        // apiErr.StatusCode, apiErr.Message, apiErr.Code
    }
}
```

## Examples

Runnable examples are in [`examples/`](examples/): quote+swap flow, providers and tokens, price lookup, swap discovery, and track. Set `SWAPKIT_API_KEY` and run e.g. `go run ./examples/quote_swap`. See [examples/README.md](examples/README.md).

## Development

- `make build` – build all packages  
- `make test` – run tests  
- `make test-cover` – coverage and HTML report  
- `make lint` – golangci-lint  
- `make fmt` – format code  
- `make tidy` – go mod tidy  

Unit tests use `httptest.Server` only; no live API calls. Integration tests against the real API require a valid API key (e.g. set `SWAPKIT_API_KEY` locally).

## Links

- [SwapKit documentation](https://docs.swapkit.dev/)
- [SwapKit API reference](https://api.swapkit.dev/docs/)
- [Register for an API key](https://dashboard.swapkit.dev/)
