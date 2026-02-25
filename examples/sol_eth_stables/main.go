// Package main demonstrates same-chain and cross-chain swaps for:
//
//   - SOL, ETH (native)
//   - USDT, USDC on Ethereum and on Solana
//
// Yes, you can do any combination of these:
//
//   - Same-chain: ETH↔ETH.USDC, ETH↔ETH.USDT, ETH.USDC↔ETH.USDT;
//     SOL↔SOL.USDC, SOL↔SOL.USDT, SOL.USDC↔SOL.USDT.
//   - Cross-chain: ETH↔SOL, ETH.USDC↔SOL.USDC, ETH.USDT↔SOL.USDT,
//     ETH.USDC↔SOL, SOL.USDC↔ETH, etc.
//
// SwapKit automatically optimizes: the API aggregates quotes from multiple providers
// and tags the best route(s) as RECOMMENDED, FASTEST, or CHEAPEST. The first route is
// usually the recommended one; use types.RouteByTag(quote, types.TagCheapest) to
// explicitly pick the cheapest route.
//
// The example uses the SDK at every level: providers, discovery (swapFrom/swapTo),
// quote for a matrix of pairs, then one full quote+swap flow.
//
// Run with SWAPKIT_API_KEY set:
//
//	SWAPKIT_API_KEY=your_key go run ./examples/sol_eth_stables
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dawitel/swapkit-sdk/client"
	"github.com/dawitel/swapkit-sdk/types"
)

// Well-known SwapKit identifiers for SOL, ETH, USDC, USDT on Ethereum and Solana.
// You can also resolve these via c.Tokens(ctx, "THORCHAIN") or c.SwapTo(ctx, "SOL.SOL").
var (
	// Native
	SOL = "SOL.SOL"
	ETH = "ETH.ETH"
	// Ethereum ERC-20 (canonical addresses)
	ETH_USDC = "ETH.USDC-0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	ETH_USDT = "ETH.USDT-0xdAC17F958D2ee523a2206206994597C13D831ec7"
	// Solana SPL (canonical mints)
	SOL_USDC = "SOL.USDC-EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	SOL_USDT = "SOL.USDT-Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"
)

func main() {
	apiKey := os.Getenv("SWAPKIT_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Set SWAPKIT_API_KEY to run this example.")
		os.Exit(1)
	}

	cfg := client.DefaultConfig(apiKey)
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	fmt.Println("=== 1. Providers (supported chains) ===")
	providers, err := c.Providers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range providers {
		hasETH := false
		hasSOL := false
		for _, ch := range p.SupportedChainIds {
			if ch == "1" || ch == "eth" {
				hasETH = true
			}
			if ch == "solana" {
				hasSOL = true
			}
		}
		if hasETH || hasSOL {
			fmt.Printf("  %s (chains: %d, tokens: %d)\n", p.Provider, len(p.SupportedChainIds), p.Count)
		}
	}

	fmt.Println("\n=== 2. Discovery: what can you swap TO with SOL.SOL? (first 8) ===")
	toList, err := c.SwapTo(ctx, SOL)
	if err != nil {
		log.Fatal(err)
	}
	for i, id := range toList {
		if i >= 8 {
			break
		}
		fmt.Printf("  %s\n", id)
	}
	fmt.Printf("  ... total %d options\n", len(toList))

	fmt.Println("\n=== 3. Quote matrix: every combination of the 6 tokens (30 pairs) ===")
	// All 6 tokens; for each sell asset, try buying every other token (5 options) = 30 directed pairs.
	tokens := []struct{ id, sellAmount string }{
		{ETH, "0.01"},
		{ETH_USDC, "100"},
		{ETH_USDT, "100"},
		{SOL, "0.1"},
		{SOL_USDC, "100"},
		{SOL_USDT, "100"},
	}
	var pairs []struct{ sell, buy, sellAmount string }
	for i, s := range tokens {
		for j := range tokens {
			if i == j {
				continue
			}
			pairs = append(pairs, struct{ sell, buy, sellAmount string }{
				sell: s.id, buy: tokens[j].id, sellAmount: s.sellAmount,
			})
		}
	}

	fmt.Printf("  Checking %d pairs (sell -> buy)...\n", len(pairs))
	for _, p := range pairs {
		quote, err := c.Quote(ctx, &types.QuoteRequest{
			SellAsset:  p.sell,
			BuyAsset:   p.buy,
			SellAmount: p.sellAmount,
			Slippage:   2,
		})
		if err != nil {
			fmt.Printf("  %s -> %s: error %v\n", p.sell, p.buy, err)
			continue
		}
		if len(quote.Routes) == 0 {
			fmt.Printf("  %s -> %s: no routes\n", p.sell, p.buy)
			continue
		}
		r := quote.Routes[0]
		outbound := r.OutboundAmount
		if outbound == "" {
			outbound = r.ExpectedBuyAmount
		}
		tags := ""
		if r.Meta != nil && len(r.Meta.Tags) > 0 {
			tags = fmt.Sprintf(" [%v]", r.Meta.Tags)
		}
		fmt.Printf("  %s -> %s: outbound %s%s\n", p.sell, p.buy, outbound, tags)
	}

	fmt.Println("\n=== 4. One full flow: Quote then Swap (placeholder addresses) ===")
	quote, err := c.Quote(ctx, &types.QuoteRequest{
		SellAsset:  ETH_USDC,
		BuyAsset:   SOL_USDC,
		SellAmount: "10",
		Slippage:   2,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Prefer the route tagged CHEAPEST when available; otherwise the first (recommended) route.
	r := types.RouteByTag(quote, types.TagCheapest)
	if r == nil {
		log.Fatal("no routes for ETH USDC -> SOL USDC")
	}
	routeID := r.RouteID
	expectedBuy := r.OutboundAmount
	if expectedBuy == "" {
		expectedBuy = r.ExpectedBuyAmount
	}
	if expectedBuy == "" {
		expectedBuy = "(see route)"
	}
	fmt.Printf("  Quote: 10 ETH.USDC -> %s SOL.USDC (routeId=%s)\n", expectedBuy, routeID)

	// swap, err := c.Swap(ctx, &types.SwapRequest{
	// 	RouteID:             routeID,
	// 	SourceAddress:       "0x0000000000000000000000000000000000000000",
	// 	DestinationAddress:  "11111111111111111111111111111111",
	// })
	// if err != nil {
	// 	fmt.Printf("  Swap (placeholders): %v\n", err)
	// 	fmt.Println("  Use real source and destination addresses to get a signable tx.")
	// 	return
	// }
	// fmt.Printf("  Swap: target=%s inbound=%s txType=%s txLen=%d\n",
	// 	swap.TargetAddress, swap.InboundAddress, swap.Meta.TxType, len(swap.Tx))
}
