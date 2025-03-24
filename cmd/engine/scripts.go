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
	app.lua.l.SetGlobal("font", app.lua.l.NewFunction(font))
	app.lua.l.SetGlobal("character", app.lua.l.NewFunction(character))
	app.lua.l.SetGlobal("narrate", app.lua.l.NewFunction(app.narrate))
	app.lua.l.SetGlobal("say", app.lua.l.NewFunction(app.say))
	app.lua.l.SetGlobal("choice", app.lua.l.NewFunction(app.choice))
	app.lua.l.SetGlobal("play_music", app.lua.l.NewFunction(app.playMusic))
	app.lua.l.SetGlobal("stop_music", app.lua.l.NewFunction(app.stopMusic))
	app.lua.l.SetGlobal("pause_music", app.lua.l.NewFunction(app.pauseMusic))
	app.lua.l.SetGlobal("resume_music", app.lua.l.NewFunction(app.resumeMusic))
	app.lua.l.SetGlobal("play_sound", app.lua.l.NewFunction(app.playSound))
	fn, err := app.lua.l.LoadFile("main.lua")
	if err != nil {
		return err
	}
	app.lua.fn = fn
	app.lua.co, _ = app.lua.l.NewThread()

	return nil
}

func character(L *lua.LState) int {
	c := L.NewTable()
	L.SetField(c, "name", L.Get(1).(lua.LString))
	L.SetField(c, "color", L.Get(2).(lua.LString))

	L.Push(c)
	return 1
}

func font(L *lua.LState) int {
	f := L.NewTable()
	L.SetField(f, "name", L.Get(1).(lua.LString))
	L.SetField(f, "path", L.Get(2).(lua.LString))

	L.Push(f)
	return 1
}

func (app *Application) narrate(L *lua.LState) int {
	// Get function arguments
	text := L.ToString(1)
	properties := L.ToTable(2)
	color, _ := hexToSDLColor(app.cfg.DefaultTextColor)
	font := app.fm.GetFont("default", 16)
	fontName := "default"
	fontPath := ""
	fontSize := 16

	if properties != nil {
		if sc, ok := properties.RawGetString("text_color").(lua.LString); ok {
			color, _ = hexToSDLColor(string(sc))
		}

		if ft, ok := properties.RawGetString("font").(*lua.LTable); ok {
			if fn, ok := ft.RawGetString("name").(lua.LString); ok {
				fontName = string(fn)
			}

			if fp, ok := ft.RawGetString("path").(lua.LString); ok {
				fontPath = string(fp)
			}
		}

		if fs, ok := properties.RawGetString("font_size").(lua.LNumber); ok {
			fontSize = int(fs)
		}
	}

	if fontPath != "" {
		if newFont, err := app.fm.LoadFont(fontName, fontPath, fontSize); err == nil {
			font = newFont
		}
	}

	// Create a text widget
	tw, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: text,
		Color: color,
		Font:  font,
	})

	// Add new widget to dialogs list
	app.dialogs.AddWidget(tw)
	app.am.Add(tw.FadeIn())

	return L.Yield(lua.LNil)
}

func (app *Application) say(L *lua.LState) int {
	// Get function arguments
	charTable := L.ToTable(1)
	char := ""
	text := L.ToString(2)
	properties := L.ToTable(3)
	charColor := sdl.Color{R: 255, G: 0, B: 0, A: 255}
	textColor, _ := hexToSDLColor(app.cfg.DefaultTextColor)
	font := app.fm.GetFont("default", 16)
	fontName := "default"
	fontPath := ""
	fontSize := 16

	if charTable != nil {
		if cn, ok := charTable.RawGetString("name").(lua.LString); ok {
			char = string(cn)
		}

		if sc, ok := charTable.RawGetString("color").(lua.LString); ok {
			charColor, _ = hexToSDLColor(string(sc))
		}
	}

	if properties != nil {
		if sc, ok := properties.RawGetString("color").(lua.LString); ok {
			textColor, _ = hexToSDLColor(string(sc))
		}

		if ft, ok := properties.RawGetString("font").(*lua.LTable); ok {
			if fn, ok := ft.RawGetString("name").(lua.LString); ok {
				fontName = string(fn)
			}

			if fp, ok := ft.RawGetString("path").(lua.LString); ok {
				fontPath = string(fp)
			}
		}

		if fs, ok := properties.RawGetString("font_size").(lua.LNumber); ok {
			fontSize = int(fs)
		}
	}

	if fontPath != "" {
		if newFont, err := app.fm.LoadFont(fontName, fontPath, fontSize); err == nil {
			font = newFont
		}
	}

	cw, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: char,
		Color: charColor,
		Font:  font,
	})

	tw, _ := gui.NewText(app.renderer, &gui.TextParams{
		Value: text,
		Color: textColor,
		Font:  font,
	})

	dw, _ := gui.NewDialog(app.renderer, &gui.DialogParams{
		Character: cw,
		Value:     tw,
	})

	app.dialogs.AddWidget(dw)
	app.am.Add(cw.FadeIn())
	app.am.Add(tw.FadeIn())

	return L.Yield(lua.LNil)
}

func (app *Application) choice(L *lua.LState) int {
	// Get function arguments
	options := L.ToTable(1)
	properties := L.ToTable(2)
	color, _ := hexToSDLColor(app.cfg.DefaultTextColor)
	font := app.fm.GetFont("default", 16)
	fontName := "default"
	fontPath := ""
	fontSize := 16

	if properties != nil {
		if sc, ok := properties.RawGetString("text_color").(lua.LString); ok {
			color, _ = hexToSDLColor(string(sc))
		}

		if ft, ok := properties.RawGetString("font").(*lua.LTable); ok {
			if fn, ok := ft.RawGetString("name").(lua.LString); ok {
				fontName = string(fn)
			}

			if fp, ok := ft.RawGetString("path").(lua.LString); ok {
				fontPath = string(fp)
			}
		}

		if fs, ok := properties.RawGetString("font_size").(lua.LNumber); ok {
			fontSize = int(fs)
		}
	}

	if fontPath != "" {
		if newFont, err := app.fm.LoadFont(fontName, fontPath, fontSize); err == nil {
			font = newFont
		}
	}

	list, _ := gui.NewList(app.renderer, &gui.ListParams{
		Spacing: 5,
	})

	for i := 1; i <= options.Len(); i++ {
		text, _ := gui.NewText(app.renderer, &gui.TextParams{
			Value: strconv.Itoa(i) + "- " + string(options.RawGetInt(i).(lua.LString)),
			Color: color,
			Font:  font,
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

func (app *Application) playMusic(L *lua.LState) int {
	path := L.ToString(1)
	loop := L.ToBool(2)

	app.aum.PlayMusic(path, loop)

	return 0
}

func (app *Application) stopMusic(L *lua.LState) int {
	app.aum.StopMusic()

	return 0
}

func (app *Application) pauseMusic(L *lua.LState) int {
	app.aum.PauseMusic()

	return 0
}

func (app *Application) resumeMusic(L *lua.LState) int {
	app.aum.ResumeMusic()

	return 0
}

func (app *Application) playSound(L *lua.LState) int {
	path := L.ToString(1)

	app.aum.PlaySound(path)

	return 0
}
