package gui

import "github.com/veandco/go-sdl2/sdl"

type Widget interface {
	Draw() (*DrawableObject, error)
	HandleEvent(sdl.Event)
	makeParent(Widget)
	getParent() Widget
	setLimit(int)
	MarkDirty()
	Destroy()
}

type DrawableObject struct {
	x       int32
	y       int32
	W       int32
	H       int32
	texture *sdl.Texture
}
