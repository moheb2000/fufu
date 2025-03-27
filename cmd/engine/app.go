package main

import (
	"time"

	"github.com/moheb2000/fufu/internal/audio"
	"github.com/moheb2000/fufu/internal/config"
	"github.com/moheb2000/fufu/internal/gui"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	lua "github.com/yuin/gopher-lua"
)

type Application struct {
	window     *sdl.Window
	renderer   *sdl.Renderer
	cfg        *config.Config
	dt         time.Duration
	fm         *gui.FontManager
	am         *gui.AnimationManager
	aum        *audio.AudioManager
	lua        *Lua
	state      string
	result     *int
	widgets    map[string]gui.Widget
	dialogs    *gui.List
	background *Background
	splash     *Splash
}

type Lua struct {
	l  *lua.LState
	co *lua.LState
	fn *lua.LFunction
}

// RunApp is responsible for initialization of SDL, running the main loop and cleanup memory
func (app *Application) RunApp() error {
	printGreeting()

	// Set default empty values for Application fields
	app.widgets = make(map[string]gui.Widget)
	app.lua = &Lua{}
	app.dt = time.Second / time.Duration(app.cfg.FPS)
	result := 0
	app.result = &result
	app.am = gui.NewAnimationManager()
	app.aum = audio.NewAudioManager(app.cfg.FPS)

	// Initialize SDL and create the main window
	err := app.initWindow()
	if err != nil {
		return err
	}

	// Create a new font manager and add a default font
	app.fm = gui.NewFontManager()
	app.fm.LoadFont("default", app.cfg.DefaultFont, 16)

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

	err = app.initScript()
	if err != nil {
		return err
	}

	// Run rhe main loop
	return app.mainLoop()
}

// cleanup free the memory at the end of the engine
func (app *Application) cleanup() {
	if app.background != nil {
		app.background.Destroy()
	}

	if app.splash != nil {
		app.splash.Destroy()
	}

	for _, widget := range app.widgets {
		widget.Destroy()
	}

	if app.renderer != nil {
		app.renderer.Destroy()
	}

	if app.window != nil {
		app.window.Destroy()
	}

	app.lua.l.Close()

	app.fm.Close()
	ttf.Quit()

	img.Quit()

	sdl.Quit()
}
