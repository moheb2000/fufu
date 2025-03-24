package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

// mainLoop runs the main loop of the engine
func (app *Application) mainLoop() error {
	// Main loop
	running := true
	for running {
		// startTime will be used at the end of the main loop to manage FPS
		startTime := time.Now()

		// Event loop
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYUP {
					if app.state == "novel" && e.Keysym.Sym == sdl.K_SPACE {
						app.lua.l.Resume(app.lua.co, app.lua.fn)
					}
				}
			}

			// Run HandleEvent function of all widgets in event loop
			for _, widget := range app.widgets {
				widget.HandleEvent(event)
			}

			if app.state == "options" && *app.result != 0 {
				app.dialogs.RemoveLastWidget()
				app.lua.l.Resume(app.lua.co, app.lua.fn, lua.LNumber(*app.result))
				*app.result = 0
				app.state = "novel"
			}
		}

		// Draw loop
		// Clear window with black color
		app.renderer.SetDrawColor(0, 0, 0, 255)
		app.renderer.Clear()

		err := app.drawLoop()
		if err != nil {
			return err
		}

		app.renderer.Present()

		// Set FPS
		remaningTime := time.Second/time.Duration(app.cfg.FPS) - time.Since(startTime)

		if remaningTime > 0 {
			app.dt = time.Second / time.Duration(app.cfg.FPS)
		} else {
			app.dt = time.Since(startTime)
		}

		if remaningTime > 0 {
			time.Sleep(remaningTime)
		}
	}

	return nil
}
