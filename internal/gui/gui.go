package gui

import (
	"log"

	"github.com/moheb2000/fufu/internal/config"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var font *ttf.Font

func init() {
	var err error

	// Initialize SDL_ttf
	if err = ttf.Init(); err != nil {
		panic(err)
	}

	cfg, err := config.Get()
	if err != nil {
		log.Fatal("Failed to open the config file:", err)
	}

	// Create a default font
	font, err = ttf.OpenFont(cfg.DefaultFont, 26*cfg.Resolution/1080)
	if err != nil {
		panic(err)
	}
}

func Close() {
	font.Close()
	ttf.Quit()
}

func Render(renderer *sdl.Renderer, do *DrawableObject) {
	renderer.Copy(do.texture, nil, &sdl.Rect{X: do.x, Y: do.y, W: do.w, H: do.h})
}
