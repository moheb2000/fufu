package main

import (
	"github.com/moheb2000/fufu/internal/gui"
	"github.com/veandco/go-sdl2/sdl"
)

func (app *Application) initDraw() error {
	bg, err := newBackground(app.renderer, "assets/background.png")
	if err != nil {
		return err
	}
	app.background = bg

	resolution, err := app.getResolution()
	if err != nil {
		return err
	}

	bgTextRect := sdl.Rect{X: 0, Y: 0, W: int32(float64(resolution.X) * 0.33), H: int32(resolution.Y)}
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

	app.am.Update(app.dt)

	app.renderer.SetLogicalSize(int32(resolution.X), int32(resolution.Y))

	return nil
}
