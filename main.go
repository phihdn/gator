package main

import (
	"fmt"
	"log"

	"github.com/phihdn/gator/internal/config"
)

func main() {
	// Read the initial configuration
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	fmt.Println("Initial config:", cfg)

	// Set the current user to "phihdn" and update the config file
	err = cfg.SetUser("phihdn")
	if err != nil {
		log.Fatalf("Error setting user: %v", err)
	}

	fmt.Println("User set successfully")

	// Read the config again to verify changes
	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading updated config: %v", err)
	}

	fmt.Printf("Updated config: %+v\n", updatedCfg)
}
