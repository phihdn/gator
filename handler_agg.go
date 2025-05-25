package main

import (
	"context"
	"fmt"
	"time"
)

// handlerAgg processes the agg command
// Usage: gator agg
func handlerAgg(s *state, cmd command) error {
	// Validate command arguments - no arguments are required
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s (no arguments required)", cmd.Name)
	}

	// Create a context with timeout to avoid hanging
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch the feed
	feedURL := "https://www.wagslane.dev/index.xml"
	fmt.Printf("Fetching feed: %s\n", feedURL)

	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}

	// Print the feed information
	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Link: %s\n", feed.Channel.Link)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)

	// Print the feed items
	fmt.Printf("\nFound %d items in feed:\n", len(feed.Channel.Items))
	for i, item := range feed.Channel.Items {
		fmt.Printf("\n--- Item %d ---\n", i+1)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Description: %s\n", item.Description)
	}

	return nil
}
