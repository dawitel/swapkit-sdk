// Package main demonstrates the quote-then-swap flow: get a quote (no addresses),
// then build a swap transaction for a chosen route (with source and destination addresses).
//
// Run with SWAPKIT_API_KEY set to use the live API:
//
//	SWAPKIT_API_KEY=your_key go run ./examples/quote_swap
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dawitel/swapkit-sdk/client"
	"github.com/dawitel/swapkit-sdk/types"
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

	// 1. Get a quote (price discovery only; no wallet addresses)
	quoteReq := &types.QuoteRequest{
		SellAsset:  "ETH.ETH",
		BuyAsset:   "BTC.BTC",
		SellAmount: "0.01",
		Slippage:   2,
	}
	quote, err := c.Quote(ctx, quoteReq)
	if err != nil {
		log.Fatal(err)
	}
	if len(quote.Routes) == 0 {
		log.Fatal("no routes returned")
	}

	fmt.Printf("Quote ID: %s\n", quote.QuoteID)
	for i, r := range quote.Routes {
		fmt.Printf("  Route %d: routeId=%s outbound=%s\n", i+1, r.RouteID, r.OutboundAmount)
		if r.Meta != nil && len(r.Meta.Tags) > 0 {
			fmt.Printf("    tags: %v\n", r.Meta.Tags)
		}
	}

	// 2. To build a swap transaction, call Swap with the chosen routeId and addresses.
	// This example does not sign or broadcast; it only shows the request shape.
	routeID := quote.Routes[0].RouteID
	swapReq := &types.SwapRequest{
		RouteID:             routeID,
		SourceAddress:       "0x0000000000000000000000000000000000000000", // replace with real sender
		DestinationAddress:  "bc1q000000000000000000000000000000000000000000", // replace with real receiver
	}
	swap, err := c.Swap(ctx, swapReq)
	if err != nil {
		// Often fails with insufficientBalance or similar when using placeholder addresses
		fmt.Printf("Swap (expected to fail with placeholders): %v\n", err)
		return
	}
	fmt.Printf("Swap ready: targetAddress=%s inboundAddress=%s txType=%s\n",
		swap.TargetAddress, swap.InboundAddress, swap.Meta.TxType)
	if len(swap.Tx) > 0 {
		fmt.Printf("Transaction payload length: %d bytes\n", len(swap.Tx))
	}
}
