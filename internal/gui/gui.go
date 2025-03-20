package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Render(renderer *sdl.Renderer, do *DrawableObject) {
	renderer.Copy(do.texture, nil, &sdl.Rect{X: do.x, Y: do.y, W: do.w, H: do.h})
}
