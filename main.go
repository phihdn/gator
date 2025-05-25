package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/phihdn/gator/internal/config"
	"github.com/phihdn/gator/internal/database"
)

// state represents the application state that is passed to command handlers
// It contains references to shared resources like configuration and database
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Load the configuration file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Connect to the database
	dbURL := cfg.DBURL
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping database: %v", err)
	}

	// Initialize database queries
	dbQueries := database.New(db)

	// Initialize application state with loaded configuration and database
	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	// Initialize the commands registry
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register all available command handlers
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	// Ensure at least one command argument is provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: gator <command> [args...]")
		return
	}

	// Parse command from arguments
	cmdName := os.Args[1]  // First argument is the command name
	cmdArgs := os.Args[2:] // Remaining arguments are passed to the command

	// Execute the requested command
	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
