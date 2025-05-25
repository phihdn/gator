package main

import (
	"context"
	"fmt"
	"os"
)

// handlerReset processes the reset command which deletes all users from the database
// This is primarily a development tool and should not be used in production
// Usage: gator reset
func handlerReset(s *state, cmd command) error {
	// Validate command arguments - no arguments are required
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s (no arguments required)", cmd.Name)
	}

	// Delete all users from the database
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		fmt.Printf("Failed to reset database: %v\n", err)
		os.Exit(1)
		return nil
	}

	// Provide user feedback
	fmt.Println("Database reset successful! All users have been deleted.")
	return nil
}
