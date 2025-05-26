# Welcome to the Blog Aggregator

We're going to build an [RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator in Go! We'll call it "Gator", you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post
- Browse your collected posts with an optional limit parameter

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## Learning Goals

- Learn how to integrate a Go application with a PostgreSQL database
- Practice using your SQL skills to query and migrate a database (using [sqlc](https://sqlc.dev/) and [goose](https://github.com/pressly/goose), two lightweight tools for typesafe SQL in Go)
- Learn how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database

## Available Commands

### User Management

- `gator register <username>` - Register a new user
- `gator login <username>` - Login as an existing user
- `gator users` - List all registered users

### Feed Management

- `gator addfeed <name> <url>` - Add a new RSS feed
- `gator feeds` - List all available feeds
- `gator follow <feed_id>` - Follow a feed
- `gator unfollow <feed_id>` - Unfollow a feed
- `gator following` - List all feeds you're following

### Content Management

- `gator browse [limit]` - View the latest posts from feeds you're following (default limit: 2)
- `gator agg <interval>` - Start the aggregator to collect posts at specified intervals (e.g., "30s", "5m")

### Utilities

- `gator reset` - Reset the database (warning: deletes all data)
