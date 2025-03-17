package gui

import "github.com/veandco/go-sdl2/sdl"

type Widget interface {
	Draw() (*DrawableObject, error)
	makeParent(Widget)
	setLimit(int)
	MarkDirty()
	Destroy()
}

type DrawableObject struct {
	x       int32
	y       int32
	w       int32
	h       int32
	texture *sdl.Texture
}
