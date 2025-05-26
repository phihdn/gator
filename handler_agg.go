package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/phihdn/gator/internal/database"
)

// handlerAgg processes the agg command
// Usage: gator agg <time_between_reqs>
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	// This line will never be reached due to the infinite loop above
	return nil
}

// scrapeFeeds fetches the next feed to process and processes it
func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

// scrapeFeed processes a single feed
func scrapeFeed(db *database.Queries, feed database.Feed) {
	// Mark the feed as fetched with the current time
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	// Fetch the feed data
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	// Process and save feed items
	savedCount := 0
	for _, item := range feedData.Channel.Items {
		// Parse the publication date
		publishedAt, err := parsePubDate(item.PubDate)
		if err != nil {
			log.Printf("Warning: could not parse pubDate for post '%s': %v", item.Title, err)
		}

		// Create a post in the database
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			// Check if it's a duplicate post error (based on URL)
			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				// Silently ignore duplicates
				continue
			}

			// Log other errors
			log.Printf("Error saving post '%s': %v", item.Title, err)
			continue
		}

		savedCount++
		fmt.Printf("Saved post: %s\n", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found, %v posts saved",
		feed.Name, len(feedData.Channel.Items), savedCount)
}
