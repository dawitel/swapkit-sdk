// Package main demonstrates swap discovery: which tokens can be swapped to/from a given asset.
//
// Run with SWAPKIT_API_KEY set:
//
//	SWAPKIT_API_KEY=your_key go run ./examples/discovery
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dawitel/swapkit-sdk/client"
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

	// Tokens you can sell to receive BTC.BTC
	buyAsset := "BTC.BTC"
	fromList, err := c.SwapFrom(ctx, buyAsset)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tokens you can SELL to get %s (%d):\n", buyAsset, len(fromList))
	for i, id := range fromList {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(fromList)-5)
			break
		}
		fmt.Printf("  %s\n", id)
	}

	// Tokens you can buy with ETH.ETH
	sellAsset := "ETH.ETH"
	toList, err := c.SwapTo(ctx, sellAsset)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTokens you can BUY with %s (%d):\n", sellAsset, len(toList))
	for i, id := range toList {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(toList)-5)
			break
		}
		fmt.Printf("  %s\n", id)
	}
}
