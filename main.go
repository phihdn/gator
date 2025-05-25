package main

import (
	"log"
	"os"

	"github.com/phihdn/gator/internal/config"
)

// state represents the application state that is passed to command handlers
// It contains references to shared resources like configuration
type state struct {
	cfg *config.Config
}

func main() {
	// Load the configuration file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Initialize application state with loaded configuration
	programState := &state{
		cfg: &cfg,
	}

	// Initialize the commands registry
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register all available command handlers
	cmds.register("login", handlerLogin)

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
