package main

import (
	"github.com/moheb2000/fufu/internal/gui"
	"github.com/veandco/go-sdl2/sdl"
)

func (app *Application) initDraw() error {
	resolution, err := app.getResolution()
	if err != nil {
		return err
	}

	bgTextRect := sdl.Rect{X: 0, Y: 0, W: int32(float64(resolution.X) * app.cfg.DialogPanel.Width), H: int32(resolution.Y)}
	if app.cfg.DialogPanel.Direction == "right" {
		bgTextRect.X = int32(resolution.X) - bgTextRect.W
	}

	wrapLength := int(app.convertLogicalToActualSizeX(int32(bgTextRect.W - bgTextRect.W/10)))

	// Create a list for showing dialogs
	list, err := gui.NewList(app.renderer, &gui.ListParams{
		Spacing:  20,
		Children: []gui.Widget{},
	})
	if err != nil {
		return err
	}
	app.dialogs = list

	// Limit the width of list
	limit, err := gui.NewLimit(app.renderer, &gui.LimitParams{
		Limit: wrapLength,
		Child: list,
	})
	if err != nil {
		return err
	}

	// Make list scrollable
	scrollable, err := gui.NewScrollableArea(app.renderer, &gui.ScrollableAreaParams{
		H:     app.convertLogicalToActualSizeY(int32(resolution.Y) * 5 / 8),
		Child: limit,
	})
	if err != nil {
		return err
	}

	// Change the position of the list
	positioned, err := gui.NewPositioned(app.renderer, &gui.PositionedParams{
		X:     app.convertLogicalToActualX(bgTextRect.X + bgTextRect.W/20),
		Y:     app.convertLogicalToActualY(bgTextRect.W / 20),
		Child: scrollable,
	})
	if err != nil {
		return err
	}

	app.widgets["dialogPanel"] = positioned

	// Create main menu
	bc, _ := hexToSDLColor(app.cfg.MainMenu.Color)
	bch, _ := hexToSDLColor(app.cfg.MainMenu.ColorHover)
	bbc, _ := hexToSDLColor(app.cfg.MainMenu.BackgroundColor)
	bbch, _ := hexToSDLColor(app.cfg.MainMenu.BackgroundColorHover)
	startValue, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: "Start",
		Color: bc,
		Font:  app.fm.GetFont("default", 16),
	})

	startButton, _ := gui.NewButton(app.renderer, &gui.ButtonParams{
		Value:                startValue,
		Color:                bc,
		ColorHover:           bch,
		BackgroundColor:      bbc,
		BackgroundColorHover: bbch,
		Width:                app.convertLogicalToActualSizeX(bgTextRect.W / 2),
		Height:               50,
	})

	startButton.OnClick(func() {
		if app.background != nil {
			app.background.Destroy()
			app.background = nil
		}

		app.state = NOVEL_STATE
		app.lua.l.Resume(app.lua.co, app.lua.fn)
		app.widgets["menu"].Destroy()
		delete(app.widgets, "menu")
	})

	quitValue, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: "Quit",
		Color: bc,
		Font:  app.fm.GetFont("default", 16),
	})

	quitButton, _ := gui.NewButton(app.renderer, &gui.ButtonParams{
		Value:                quitValue,
		Color:                bc,
		ColorHover:           bch,
		BackgroundColor:      bbc,
		BackgroundColorHover: bbch,
		Width:                app.convertLogicalToActualSizeX(bgTextRect.W / 2),
		Height:               50,
	})

	// Add callback function to quit button that emit sdl.QUIT event
	quitButton.OnClick(func() {
		sdl.PushEvent(&sdl.QuitEvent{Type: sdl.QUIT, Timestamp: sdl.GetTicks()})
	})

	menuList, _ := gui.NewList(app.renderer, &gui.ListParams{
		Spacing: 20,
		Children: []gui.Widget{
			startButton,
			quitButton,
		},
	})

	menu, _ := gui.NewPositioned(app.renderer, &gui.PositionedParams{
		X:     app.convertLogicalToActualX(bgTextRect.X + bgTextRect.W/2 - bgTextRect.W/4),
		Y:     app.convertLogicalToActualY(bgTextRect.Y),
		Child: menuList,
	})

	app.widgets["menu"] = menu

	// Create main menu background
	if app.cfg.MainMenu.Background != "" {
		mainMenuBackground, err := newBackground(app.renderer, &BackgroundParams{
			Path: app.cfg.MainMenu.Background,
			Origin: &Origin{
				X: "left",
				Y: "left",
			},
			App: app,
		})
		if err != nil {
			return err
		}

		app.background = mainMenuBackground
	}
	return nil
}

func (app *Application) drawLoop() error {
	// Get resolution for calculating position and size of objects
	resolution, err := app.getResolution()
	if err != nil {
		return err
	}

	// Define dialog panel rectangle
	bgTextRect := sdl.Rect{X: 0, Y: 0, W: int32(float64(resolution.X) * app.cfg.DialogPanel.Width), H: int32(resolution.Y)}
	if app.cfg.DialogPanel.Direction == "right" {
		bgTextRect.X = int32(resolution.X) - bgTextRect.W
	}

	if app.background != nil {
		app.background.draw()
	}

	if app.background == nil && app.state == MENU_STATE {
		bgColor, _ := hexToSDLColor(app.cfg.MainMenu.BackgroundColor)
		app.renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, 255)
		app.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(resolution.X), H: int32(resolution.Y)})
	}

	// Draw dialogs' background
	dpc, err := hexToSDLColor(app.cfg.DialogPanel.Color)
	if err != nil {
		app.renderer.SetDrawColor(20, 20, 20, 255)
	} else {
		app.renderer.SetDrawColor(dpc.R, dpc.G, dpc.B, 255)
	}
	app.renderer.FillRect(&bgTextRect)

	app.renderer.SetLogicalSize(0, 0)

	// This is because the window size is changed after event handling in main loop
	if pos, ok := app.widgets["menu"].(*gui.Positioned); ok {
		pdo, _ := pos.Draw()
		pos.SetPosition(app.convertLogicalToActualX(bgTextRect.X+bgTextRect.W/2-bgTextRect.W/4), app.convertLogicalToActualY(bgTextRect.Y+bgTextRect.H/2)-pdo.H/2)
	}

	for _, w := range app.widgets {
		w1, _ := w.Draw()
		gui.Render(app.renderer, w1)
	}

	app.am.Update(app.dt)

	app.renderer.SetLogicalSize(int32(resolution.X), int32(resolution.Y))

	// Draw splash
	if app.splash != nil {
		done, err := app.splash.draw()
		if err != nil {
			return err
		}

		if done {
			app.splash.Destroy()
			app.splash = nil
			app.state = NOVEL_STATE
		}
	}

	return nil
}
