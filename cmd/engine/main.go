// Fufu is a visual novel engine written in golang.
package main

import (
	"log"
	"time"

	"github.com/moheb2000/fufu/internal/config"
	"github.com/moheb2000/fufu/internal/gui"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

type Application struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	cfg        *config.Config
	widgets    map[string]gui.Widget
	background *Background
}

// RunApp is responsible for initialization of SDL, running the main loop and cleanup memory
func (app *Application) RunApp() error {
	printGreeting()

	app.widgets = make(map[string]gui.Widget)

	// Initialize SDL and create the main window
	err := app.initWindow()
	if err != nil {
		return err
	}

	// Initialize the main renderer
	err = app.initRenderer()
	if err != nil {
		return err
	}

	// Draw at the startup
	err = app.initDraw()
	if err != nil {
		return err
	}

	// Free memory at the end of the RunAPpp function
	defer app.cleanup()

	// Start lua VM and compiler
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("say", L.NewFunction(app.say))
	if err := L.DoFile("main.lua"); err != nil {
		log.Fatal(err)
	}

	// Run rhe main loop
	return app.mainLoop()
}

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

func (app *Application) initDraw() error {
	bg, err := newBackground(app.renderer, "assets/background.png")
	if err != nil {
		return err
	}
	app.background = bg

	return nil
}

func (app *Application) drawLoop() error {
	// Get resolution for calculating position and size of objects
	resolution, err := app.getResolution()
	if err != nil {
		return err
	}

	app.background.draw(app.renderer)

	// Define dialog panel rectangle
	bgTextRect := sdl.Rect{X: 0, Y: 0, W: int32(float64(resolution.X) * 0.33), H: int32(resolution.Y)}
	if app.cfg.DialogPanel.Direction == "right" {
		bgTextRect.X = int32(resolution.X) - bgTextRect.W
	}

	// Draw dialogs' background
	app.renderer.SetDrawColor(20, 20, 20, 255)
	app.renderer.FillRect(&bgTextRect)

	app.renderer.SetLogicalSize(0, 0)

	for _, w := range app.widgets {
		w1, _ := w.Draw()
		gui.Render(app.renderer, w1)
	}

	app.renderer.SetLogicalSize(int32(resolution.X), int32(resolution.Y))

	return nil
}

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
					if e.Keysym.Sym == sdl.K_SPACE {
						// TODO: Handle spaces to show the next dialog here
					}
				}
			case *sdl.MouseWheelEvent:
				// TODO: Handle mouse wheel event for scrolling in widgets. Perhaps using HandleEvent for each widget is better than handling scrolling here
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
			time.Sleep(remaningTime)
		}
	}

	return nil
}

// cleanup free the memory at the end of the engine
func (app *Application) cleanup() {
	app.background.Destroy()

	for _, widget := range app.widgets {
		widget.Destroy()
	}

	if app.renderer != nil {
		app.renderer.Destroy()
	}

	if app.window != nil {
		app.window.Destroy()
	}

	sdl.Quit()
}

func (app *Application) say(L *lua.LState) int {
	goodbye, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: L.ToString(1),
		Color: sdl.Color{R: 255, G: 255, B: 255, A: 255},
	})
	app.widgets["goodbye"] = goodbye

	return 0
}

// main is the starting point of engine.
func main() {
	// Get engine configs from "config.json" file
	cfg, err := config.Get()
	if err != nil {
		log.Fatal("Failed to open the config file:", err)
	}

	app := Application{
		cfg: cfg,
	}

	err = app.RunApp()
	if err != nil {
		log.Println("[ERROR]", err)
	}
}
