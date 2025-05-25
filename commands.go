package main

import "errors"

// command represents a CLI command with a name and arguments
// Name is the command identifier (e.g., "login")
// Args contains any additional arguments passed to the command
type command struct {
	Name string
	Args []string
}

// commands stores all registered command handlers
// It uses a map to associate command names with their handler functions
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// register adds a new command handler to the commands registry
// name: the command identifier (e.g., "login")
// f: the handler function that will be called when this command is executed
func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

// run executes a command if it exists in the registry
// It finds the appropriate handler based on cmd.Name and executes it
// Returns an error if the command doesn't exist or if the handler returns an error
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
