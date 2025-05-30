package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/phihdn/gator/internal/database"
)

// handlerAddFeed processes the addfeed command, which adds a new feed to the database
// It expects two arguments: the name of the feed and its URL
// Usage: gator addfeed <name> <url>
func handlerAddFeed(s *state, cmd command, user database.User) error {
	// Validate command arguments - exactly two args required (name and url)
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

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
		// Check if this is a duplicate feed URL error
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" && strings.Contains(pgErr.Message, "feeds_url_key") {
			// If the feed already exists, try to follow it
			existingFeed, err := s.db.GetFeedByURL(context.Background(), url)
			if err != nil {
				return fmt.Errorf("error retrieving existing feed: %w", err)
			}

			// Try to create a feed follow for the existing feed
			feedFollowID := uuid.New()
			_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
				ID:        feedFollowID,
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    user.ID,
				FeedID:    existingFeed.ID,
			})

			if err != nil {
				// Check if user is already following this feed
				if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" && strings.Contains(pgErr.Message, "feed_follows_user_id_feed_id_key") {
					fmt.Printf("Feed with URL '%s' already exists and you are already following it.\n", url)
					return nil
				}
				return fmt.Errorf("couldn't follow existing feed: %w", err)
			}

			fmt.Printf("Feed with URL '%s' already exists. You are now following it.\n", url)
			return nil
		}
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	// Automatically follow the feed after creating it
	feedFollowID := uuid.New()
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        feedFollowID,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("couldn't follow the feed: %w", err)
	}

	// Display the feed information
	fmt.Println("Feed added successfully:")
	fmt.Printf("ID: %s\n", feed.ID)
	fmt.Printf("Name: %s\n", feed.Name)
	fmt.Printf("URL: %s\n", feed.Url)
	fmt.Printf("Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))
	fmt.Printf("You are now following this feed.\n")

	return nil
}

// handlerFeeds processes the feeds command, which lists all feeds in the database
// It takes no arguments and prints all feeds along with their owner's name
// Usage: gator feeds
func handlerFeeds(s *state, cmd command, user database.User) error {
	// Validate command arguments - no args expected
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s (takes no arguments)", cmd.Name)
	}

	// Get all feeds with associated user information
	feeds, err := s.db.GetAllFeedsWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	// Check if there are feeds to display
	if len(feeds) == 0 {
		fmt.Println("No feeds found in the database")
		return nil
	}

	// Display feed information
	fmt.Printf("Found %d feeds:\n\n", len(feeds))
	for i, feed := range feeds {
		fmt.Printf("Feed #%d:\n", i+1)
		fmt.Printf("  Name: %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  Created By: %s\n", feed.UserName)
		fmt.Printf("  Added On: %s\n\n", feed.CreatedAt.Format(time.RFC3339))
	}

	return nil
}

// handlerFollowFeed processes the follow command, which allows a user to follow a feed
// It takes a single URL argument and creates a feed follow record for the current user
// Usage: gator follow <url>
func handlerFollowFeed(s *state, cmd command, user database.User) error {
	// Validate command arguments - exactly one arg required (url)
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	// Find the feed by URL
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no feed found with URL '%s'", url)
		}
		return fmt.Errorf("error finding feed: %w", err)
	}

	// Create a new feed follow record
	now := time.Now().UTC()
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		// Check if this is a duplicate error (user already following this feed)
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && strings.Contains(pgErr.Message, "feed_follows_user_id_feed_id_key") {
				fmt.Printf("You are already following the feed '%s'\n", feed.Name)
				return nil
			}
		}
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	// Print confirmation message
	fmt.Printf("You are now following the feed '%s'\n", feedFollow.FeedName)

	return nil
}

// handlerFollowing processes the following command, which lists all feeds a user is following
// It takes no arguments and displays all feeds the current user is following
// Usage: gator following
func handlerFollowing(s *state, cmd command, user database.User) error {
	// Validate command arguments - no args expected
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s (takes no arguments)", cmd.Name)
	}

	// Get all feed follows for the current user
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	// Check if there are any feed follows to display
	if len(feedFollows) == 0 {
		fmt.Printf("User '%s' is not following any feeds\n", user.Name)
		return nil
	}

	// Display feed follow information
	fmt.Printf("User '%s' is following %d feeds:\n\n", user.Name, len(feedFollows))
	for i, followedFeed := range feedFollows {
		fmt.Printf("Feed #%d: %s\n", i+1, followedFeed.FeedName)
		fmt.Printf("  Followed on: %s\n\n", followedFeed.CreatedAt.Format(time.RFC3339))
	}

	return nil
}

// handlerUnfollowFeed processes the unfollow command, which allows a user to stop following a feed
// It takes a single URL argument and removes the feed follow record for the current user
// Usage: gator unfollow <url>
func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	// Validate command arguments - exactly one arg required (url)
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	// Find the feed by URL first, to provide better error messages
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no feed found with URL '%s'", url)
		}
		return fmt.Errorf("error finding feed: %w", err)
	}

	// Delete the feed follow
	err = s.db.DeleteFeedFollowByUserAndFeedURL(context.Background(), database.DeleteFeedFollowByUserAndFeedURLParams{
		UserID: user.ID,
		Url:    url,
	})

	if err != nil {
		// Check if it was a "no rows affected" situation
		if err == sql.ErrNoRows {
			fmt.Printf("You are not currently following the feed with URL '%s'\n", url)
			return nil
		}
		return fmt.Errorf("error unfollowing feed: %w", err)
	}

	// Print confirmation message
	fmt.Printf("You have successfully unfollowed the feed '%s'\n", feed.Name)

	return nil
}
