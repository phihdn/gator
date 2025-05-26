package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

// RSSFeed represents a parsed RSS feed
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem represents a single item in an RSS feed
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// parsePubDate attempts to parse the pubDate string from an RSS feed
// It tries multiple common time formats used in RSS feeds
func parsePubDate(pubDate string) (sql.NullTime, error) {
	if pubDate == "" {
		return sql.NullTime{Valid: false}, nil
	}

	// List of common time formats used in RSS feeds
	timeFormats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"02 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 MST",
	}

	for _, format := range timeFormats {
		parsedTime, err := time.Parse(format, pubDate)
		if err == nil {
			return sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}, nil
		}
	}

	// If we couldn't parse the time with any format
	return sql.NullTime{Valid: false}, fmt.Errorf("could not parse date: %s", pubDate)
}

// fetchFeed fetches and parses an RSS feed from the given URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new HTTP request with context
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set the User-Agent header to identify our client
	request.Header.Set("User-Agent", "gator")

	// Create an HTTP client and execute the request
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer response.Body.Close()

	// Check if the response was successful
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the XML content
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("error parsing feed XML: %w", err)
	}

	// Unescape HTML entities in the feed title and description
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Unescape HTML entities in each item's title and description
	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
		feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
	}

	return &feed, nil
}
