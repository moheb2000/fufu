package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// initWindow initializes SDL and create the main window
func (app *Application) initWindow() error {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return err
	}

	// Initialize SDL_img
	if err := img.Init(img.INIT_PNG | img.INIT_JPG); err != nil {
		return err
	}

	// Initialize SDL_ttf
	if err := ttf.Init(); err != nil {
		return err
	}

	// Create main window
	window, err := sdl.CreateWindow(app.cfg.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1024, 720, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return err
	}
	app.window = window

	// Make window fullscreen based on user config
	if app.cfg.FullScreen {
		app.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	}

	return nil
}
