// Helper functions for the engine
package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Resolution struct {
	X int
	Y int
}

// printGreeting prints some information about engine and libraries that the starting point of engine
func printGreeting() {
	// Engine name and repository address
	fmt.Printf("Fufu Visual Novel Engine - https://github.com/moheb2000/fufu\n")

	// SDL library version
	sdlVersion := sdl.Version{}
	sdl.GetVersion(&sdlVersion)
	fmt.Printf("SDL Version: %d.%d.%d\n", sdlVersion.Major, sdlVersion.Minor, sdlVersion.Patch)
}

// getResolution returns the user resolution based on config.go
func (app *Application) getResolution() (*Resolution, error) {
	r := Resolution{}
	// Check game resolution based on user config
	r.Y = app.cfg.Resolution
	switch r.Y {
	case 720:
		r.X = 1280
	case 1080:
		r.X = 1920
	case 2160:
		r.X = 3840
	case 4320:
		r.X = 7680
	default:
		return nil, fmt.Errorf("%d resolution does not supported! Supported resolutions: 720(HD), 1080(Full HD), 2160(4K), 4320(8K)", r.Y)
	}

	return &r, nil
}

// getActualLogicalSize returns the actual width and height of the renderer area when logical size set for renderer
func (app *Application) getActualLogicalSize() (int32, int32, error) {
	// Get resolution for renderer
	resolution, err := app.getResolution()
	if err != nil {
		return 0, 0, err
	}

	// Get actual window width and height. (GetOutputSize is more preferable than GetWindowSize in SDL because it handles high DPI monitors better)
	winX, winY, err := app.renderer.GetOutputSize()
	if err != nil {
		return 0, 0, err
	}

	// Check the scaling between window size and resolution
	scaleX := float64(winX) / float64(resolution.X)
	scaleY := float64(winY) / float64(resolution.Y)

	// Compute the actual size of rendering area based on the comparison of scales
	var actualX int32
	var actualY int32
	if scaleX <= scaleY {
		actualX = winX
		actualY = int32(float64(resolution.Y) * scaleX)
	} else {
		actualY = winY
		actualX = int32(float64(resolution.X) * scaleY)
	}

	return actualX, actualY, nil
}

// Because based on whether we want the actual position based on the logical position is dependent if the position is on x coordinate or y, we define two functions to cover both. For simplisity in using these functions, we don't handle errors in them, becuase it's unlikely to happen here

// convertLogicalToActualX gets a logical position on x coordinate and convert it to actual x on the same coordinate
func (app *Application) convertLogicalToActualX(lp int32) int32 {
	// Get necassery values for calculating
	resolution, _ := app.getResolution()
	winX, _, _ := app.renderer.GetOutputSize()
	actualX, _, _ := app.getActualLogicalSize()

	// (winX-actualX)/2 is because the rendering area is centered by default and we need to add the black area in left side of it for getting the actual size
	return lp*actualX/int32(resolution.X) + (winX-actualX)/2
}

// convertLogicalToActualY gets a logical position on y coordinate and convert it to actual y on the same coordinate
func (app *Application) convertLogicalToActualY(lp int32) int32 {
	// Get necassery values for calculating
	resolution, _ := app.getResolution()
	_, winY, _ := app.renderer.GetOutputSize()
	_, actualY, _ := app.getActualLogicalSize()

	// (winX-actualY)/2 is because the rendering area is centered by default and we need to add the black area in top side of it for getting the actual size
	return lp*actualY/int32(resolution.Y) + (winY-actualY)/2
}

func (app *Application) convertLogicalToActualSizeX(ls int32) int32 {
	winX, _, _ := app.renderer.GetOutputSize()
	actualX, _, _ := app.getActualLogicalSize()

	return app.convertLogicalToActualX(ls) - (winX-actualX)/2
}

func (app *Application) convertLogicalToActualSizeY(ls int32) int32 {
	_, winY, _ := app.renderer.GetOutputSize()
	_, actualY, _ := app.getActualLogicalSize()

	return app.convertLogicalToActualY(ls) - (winY-actualY)/2
}
