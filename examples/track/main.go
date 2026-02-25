// Package main demonstrates tracking a swap by transaction hash and chain ID.
//
// Run with SWAPKIT_API_KEY set. Pass hash and chainId as args, or use placeholders:
//
//	SWAPKIT_API_KEY=your_key go run ./examples/track [hash] [chainId]
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

	hash := "0x1890aba1c0b25126892af2ab09f5c1bba75adefc47918a96ea498764ab643ce9"
	chainID := "1"
	if len(os.Args) >= 3 {
		hash = os.Args[1]
		chainID = os.Args[2]
	}

	cfg := client.DefaultConfig(apiKey)
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	req := &types.TrackRequest{Hash: hash, ChainID: chainID}
	resp, err := c.Track(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %s (tracking: %s)\n", resp.Status, resp.TrackingStatus)
	fmt.Printf("From: %s %s -> To: %s %s\n", resp.FromAsset, resp.FromAmount, resp.ToAsset, resp.ToAmount)
	fmt.Printf("Chain: %s Hash: %s Block: %d\n", resp.ChainID, resp.Hash, resp.Block)
	if resp.Meta != nil && resp.Meta.Provider != "" {
		fmt.Printf("Provider: %s\n", resp.Meta.Provider)
	}
	if len(resp.Legs) > 0 {
		fmt.Printf("Legs: %d\n", len(resp.Legs))
	}
}
