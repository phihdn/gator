# Welcome to the Blog Aggregator

We're going to build an [RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator in Go! We'll call it "Gator", you know, because aggreGATOR 🐊. Anyhow, it's a CLI tool that allows users to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post
- Browse your collected posts with an optional limit parameter

RSS feeds are a way for websites to publish updates to their content. You can use this project to keep up with your favorite blogs, news sites, podcasts, and more!

## Prerequisites

To run Gator, you'll need:

- [Go](https://golang.org/doc/install) (version 1.16 or later)
- [PostgreSQL](https://www.postgresql.org/download/) database server (version 12 or later)

## Installation

You can install the Gator CLI directly from the source code:

```bash
# Clone the repository (if you haven't already)
git clone https://github.com/yourusername/gator.git
cd gator

# Install the CLI
go install .
```

Alternatively, if the package is published:

```bash
go install github.com/yourusername/gator@latest
```

## Configuration

Gator requires a PostgreSQL database connection. The application will look for a configuration file at `~/.gatorconfig.json` in your home directory with the following structure:

```json
{
  "db_url": "postgresql://username:password@localhost:5432/gator",
  "current_user_name": ""
}
```

Replace `username`, `password`, and other database connection details with your own PostgreSQL configuration. The `current_user_name` will be automatically updated when you login.

## Database Setup

Before using Gator, you need to set up the database schema using the SQL migrations:

1. Install goose (database migration tool):

   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. Create a PostgreSQL database named `gator`:

   ```bash
   createdb gator
   ```

3. Run the migrations:

   ```bash
   cd /path/to/gator
   goose -dir sql/schema postgres "postgresql://username:password@localhost:5432/gator" up
   ```

   Replace the connection string with your PostgreSQL credentials.

## Getting Started

1. Make sure PostgreSQL is running
2. Register a new user:

   ```bash
   gator register yourusername
   ```

3. Login as the user:

   ```bash
   gator login yourusername
   ```

4. Add an RSS feed to follow:

   ```bash
   gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
   ```

5. List available feeds:

   ```bash
   gator feeds
   ```

6. Follow a feed (using the feed_id from the list):

   ```bash
   gator follow 1
   ```

7. Start collecting posts:

   ```bash
   gator agg 5m
   ```

8. In another terminal, browse your collected posts:

   ```bash
   gator browse 5
   ```

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

## Extending the Project

here are some ideas:

- [ ] Add sorting and filtering options to the browse command
- [ ] Add pagination to the browse command
- [ ] Add concurrency to the agg command so that it can fetch more frequently
- [ ] Add a search command that allows for fuzzy searching of posts
- [ ] Add bookmarking or liking posts
- [ ] Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
- [ ] Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
- [ ] Write a service manager that keeps the agg command running in the background and restarts it if it crashes