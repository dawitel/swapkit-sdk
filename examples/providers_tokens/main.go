// Package main lists swap providers and then fetches supported tokens for one provider.
//
// Run with SWAPKIT_API_KEY set:
//
//	SWAPKIT_API_KEY=your_key go run ./examples/providers_tokens
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

	// List all providers and their supported chains
	providers, err := c.Providers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Providers (%d):\n", len(providers))
	for _, p := range providers {
		fmt.Printf("  %s (provider=%s) count=%d chains=%v\n",
			p.Name, p.Provider, p.Count, p.SupportedChainIds)
	}

	if len(providers) == 0 {
		return
	}

	// Fetch tokens for the first provider
	providerName := providers[0].Provider
	tokens, err := c.Tokens(ctx, providerName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTokens for %s (%d):\n", tokens.Provider, tokens.Count)
	for i, t := range tokens.Tokens {
		if i >= 5 {
			fmt.Printf("  ... and %d more\n", len(tokens.Tokens)-5)
			break
		}
		fmt.Printf("  %s %s (identifier=%s)\n", t.Chain, t.Ticker, t.Identifier)
	}
}
