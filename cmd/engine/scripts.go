package main

import (
	"strconv"

	"github.com/moheb2000/fufu/internal/gui"
	"github.com/veandco/go-sdl2/sdl"
	lua "github.com/yuin/gopher-lua"
)

func (app *Application) initScript() error {
	// Start lua VM and compiler
	app.lua.l = lua.NewState()
	app.lua.l.SetGlobal("say", app.lua.l.NewFunction(app.say))
	app.lua.l.SetGlobal("options", app.lua.l.NewFunction(app.options))
	fn, err := app.lua.l.LoadFile("main.lua")
	if err != nil {
		return err
	}
	app.lua.fn = fn
	app.lua.co, _ = app.lua.l.NewThread()

	return nil
}

func (app *Application) say(L *lua.LState) int {
	text, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: L.ToString(1),
		Color: sdl.Color{R: 255, G: 255, B: 255, A: 255},
		Font:  app.fm.GetFont("default", 16),
	})
	app.dialogs.AddWidget(text)
	app.am.Add(text.FadeIn())

	return L.Yield(lua.LNil)
}

func (app *Application) options(L *lua.LState) int {
	nArgs := L.GetTop()

	list, _ := gui.NewList(app.renderer, &gui.ListParams{
		// Children: []gui.Widget{},
		Spacing: 5,
	})

	for i := 1; i <= nArgs; i++ {
		text, _ := gui.NewText(app.renderer, &gui.TextParams{
			Value: strconv.Itoa(i) + "- " + L.ToString(i),
			Color: sdl.Color{R: 255, G: 255, B: 0, A: 255},
			Font:  app.fm.GetFont("default", 16),
		})

		list.AddWidget(text)
		app.am.Add(text.FadeIn())

		app.state = "options"
	}

	ops, _ := gui.NewOptions(app.renderer, &gui.OptionsParams{
		Options: list,
		Result:  app.result,
	})

	app.dialogs.AddWidget(ops)

	L.Push(lua.LNil)

	return L.Yield(lua.LNumber(1))
}
