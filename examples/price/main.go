// Package main fetches token prices for one or more identifiers, with optional metadata.
//
// Run with SWAPKIT_API_KEY set:
//
//	SWAPKIT_API_KEY=your_key go run ./examples/price
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

	req := &types.PriceRequest{
		Tokens: []types.PriceTokenInput{
			{Identifier: "ETH.ETH"},
			{Identifier: "BTC.BTC"},
			{Identifier: "SOL.SOL"},
		},
		Metadata: true,
	}
	prices, err := c.Price(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Prices:")
	for _, p := range prices {
		fmt.Printf("  %s: $%.2f (timestamp %d)\n", p.Identifier, p.PriceUSD, p.Timestamp)
		if p.Cg != nil {
			fmt.Printf("    CoinGecko: %s (id=%s)\n", p.Cg.Name, p.Cg.ID)
		}
	}
}
