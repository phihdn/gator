package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/phihdn/gator/internal/database"
)

// handlerLogin processes the login command which sets the current user in the config
// It validates that exactly one argument (username) is provided and that the user exists in the database
// Usage: gator login <username>
func handlerLogin(s *state, cmd command) error {
	// Validate command arguments - exactly one username is required
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	name := cmd.Args[0]

	// Check if the user exists in the database
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("User '%s' does not exist, please register first.\n", name)
			os.Exit(1)
		}
		return fmt.Errorf("couldn't find user: %w", err)
	}

	// Update the user in configuration
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	// Provide user feedback
	fmt.Println("User switched successfully!")
	return nil
}

// handlerRegister processes the register command which creates a new user in the database
// It validates that exactly one argument (username) is provided
// Usage: gator register <username>
func handlerRegister(s *state, cmd command) error {
	// Validate command arguments - exactly one username is required
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	name := cmd.Args[0]

	// Create a new user in the database
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	// Check if there was an error creating the user
	if err != nil {
		// Check for a unique constraint violation, which would indicate the user already exists
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_name_key\"" {
			fmt.Printf("User '%s' already exists!\n", name)
			os.Exit(1)
		}
		return fmt.Errorf("couldn't create user: %w", err)
	}

	// Update the user in configuration
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	// Provide user feedback and log the user's data
	fmt.Printf("User '%s' created successfully!\n", name)
	log.Printf("User created: %+v\n", user)
	return nil
}
