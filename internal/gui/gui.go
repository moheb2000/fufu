package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

func Render(renderer *sdl.Renderer, do *DrawableObject) {
	renderer.Copy(do.texture, nil, &sdl.Rect{X: do.x, Y: do.y, W: do.w, H: do.h})
}

// TODO: I don't know how blend mode works, so I don't know is this a correct approach or not but for now it increases the text quallity so I use it
var BLENDMOD_ONE = sdl.ComposeCustomBlendMode(
	sdl.BLENDFACTOR_ONE, sdl.BLENDFACTOR_ONE,
	sdl.BLENDOPERATION_ADD,
	sdl.BLENDFACTOR_ONE, sdl.BLENDFACTOR_ONE,
	sdl.BLENDOPERATION_ADD,
)
