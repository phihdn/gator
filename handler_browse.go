package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/phihdn/gator/internal/database"
)

// handlerBrowse processes the browse command
// Usage: gator browse [limit]
func handlerBrowse(s *state, cmd command, user database.User) error {
	// Default limit to 2 if not provided
	limit := int32(2)

	// Parse limit from args if provided
	if len(cmd.Args) >= 1 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(parsedLimit)
	}

	// Get posts for the user
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}

	// Print the posts
	fmt.Printf("Found %d posts for %s:\n\n", len(posts), user.Name)
	for _, post := range posts {
		var publishedAt string
		if post.PublishedAt.Valid {
			publishedAt = post.PublishedAt.Time.Format("Jan 02, 2006")
		} else {
			publishedAt = "unknown date"
		}

		fmt.Printf("Feed: %s\n", post.FeedName)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Published: %s\n", publishedAt)
		fmt.Printf("URL: %s\n", post.Url)

		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}

		fmt.Println("--------------------")
	}

	if len(posts) == 0 {
		fmt.Println("No posts found. Try following some feeds first!")
	}

	return nil
}
