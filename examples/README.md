# Examples

Runnable examples for the SwapKit Go SDK. Each requires `SWAPKIT_API_KEY` to be set (get one at [dashboard.swapkit.dev](https://dashboard.swapkit.dev/)).

| Example | Description |
|---------|-------------|
| **quote_swap** | Quote-then-swap flow: get a quote (no addresses), then build a swap for a chosen route. |
| **providers_tokens** | List all providers and their supported chains; fetch tokens for one provider. |
| **price** | Fetch USD prices for ETH, BTC, SOL with optional CoinGecko metadata. |
| **discovery** | Swap discovery: tokens you can sell to get BTC; tokens you can buy with ETH. |
| **track** | Track a swap by transaction hash and chain ID (optional args: `hash` `chainId`). |
| **sol_eth_stables** | SOL, ETH, USDC/USDT on ETH and SOL: same-chain and cross-chain quote matrix and one full quote+swap flow. |

Run any example:

```bash
export SWAPKIT_API_KEY=your_key
go run ./examples/quote_swap
go run ./examples/providers_tokens
go run ./examples/price
go run ./examples/discovery
go run ./examples/track
go run ./examples/sol_eth_stables
# With custom hash/chain for track:
go run ./examples/track 0xYourTxHash 1
```
