package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/phihdn/gator/internal/database"
)

// middlewareLoggedIn is middleware that ensures a user is logged in
// It takes a handler function that requires a logged-in user and returns
// a normal handler function that can be registered with the command system
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Get current user from config
		currentUserName := s.cfg.CurrentUserName
		if currentUserName == "" {
			fmt.Println("No user is logged in. Please login first.")
			os.Exit(1)
		}

		// Get current user from database
		user, err := s.db.GetUser(context.Background(), currentUserName)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Printf("User '%s' does not exist, please register first.\n", currentUserName)
				os.Exit(1)
			}
			return fmt.Errorf("couldn't find user: %w", err)
		}

		// Call the original handler with the user
		return handler(s, cmd, user)
	}
}
