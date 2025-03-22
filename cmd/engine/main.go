// Fufu is a visual novel engine written in golang.
package main

import (
	"log"

	"github.com/moheb2000/fufu/internal/config"
)

// main is the starting point of engine.
func main() {
	// Get engine configs from "config.json" file
	cfg, err := config.Get()
	if err != nil {
		log.Fatal("Failed to open the config file:", err)
	}

	app := Application{
		cfg:   cfg,
		state: "novel",
	}

	err = app.RunApp()
	if err != nil {
		log.Println("[ERROR]", err)
	}
}
