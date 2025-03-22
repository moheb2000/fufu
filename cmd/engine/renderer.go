package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

// initRenderer initializes the main renderer
func (app *Application) initRenderer() error {
	// Create main renderer
	renderer, err := sdl.CreateRenderer(app.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal("Failed to create the main renderer:", err)
	}
	app.renderer = renderer

	// Check game resolution based on user config
	resolution, err := app.getResolution()
	if err != nil {
		return err
	}

	// SetLogicalSize ensures game runs in a consistent aspect ratio while window size changes
	app.renderer.SetLogicalSize(int32(resolution.X), int32(resolution.Y))

	return nil
}
