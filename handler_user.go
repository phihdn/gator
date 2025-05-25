package main

import (
	"fmt"
)

// handlerLogin processes the login command which sets the current user in the config
// It validates that exactly one argument (username) is provided
// Usage: gator login <username>
func handlerLogin(s *state, cmd command) error {
	// Validate command arguments - exactly one username is required
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	name := cmd.Args[0]

	// Update the user in configuration
	err := s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	// Provide user feedback
	fmt.Println("User switched successfully!")
	return nil
}
