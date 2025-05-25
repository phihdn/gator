package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/phihdn/gator/internal/database"
)

// handlerAddFeed processes the addfeed command, which adds a new feed to the database
// It expects two arguments: the name of the feed and its URL
// Before adding the feed, it checks if the current user exists in the database
// Usage: gator addfeed <name> <url>
func handlerAddFeed(s *state, cmd command) error {
	// Validate command arguments - exactly two args required (name and url)
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	// Get current user from config
	currentUserName := s.cfg.CurrentUserName
	if currentUserName == "" {
		fmt.Println("No user is logged in. Please login first.")
		os.Exit(1)
	}

	// Get current user from database
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("User '%s' does not exist, please register first.\n", currentUserName)
			os.Exit(1)
		}
		return fmt.Errorf("couldn't find user: %w", err)
	}

	// Create a new feed record
	now := time.Now().UTC()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	// Display the feed information
	fmt.Println("Feed added successfully:")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))

	return nil
}
