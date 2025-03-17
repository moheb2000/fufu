package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Limit struct {
	parent         Widget
	dirty          bool
	renderer       *sdl.Renderer
	limitParams    *LimitParams
	drawableObject *DrawableObject
}

type LimitParams struct {
	Limit int
	Child Widget
}

func NewLimit(renderer *sdl.Renderer, p *LimitParams) (*Limit, error) {
	l := Limit{
		renderer:       renderer,
		limitParams:    p,
		dirty:          true,
		drawableObject: &DrawableObject{},
	}

	p.Child.makeParent(&l)

	l.setLimit(l.limitParams.Limit)

	err := l.updateTexture()
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *Limit) updateTexture() error {
	if !l.dirty {
		return nil
	}

	do, err := l.limitParams.Child.Draw()
	if err != nil {
		return err
	}

	l.drawableObject = do

	l.dirty = false

	return nil
}

func (l *Limit) Draw() (*DrawableObject, error) {
	err := l.updateTexture()
	if err != nil {
		return nil, err
	}

	return l.drawableObject, nil
}

func (l *Limit) makeParent(parent Widget) {
	l.parent = parent
}

func (l *Limit) setLimit(limit int) {
	l.limitParams.Child.setLimit(limit)
}

func (l *Limit) MarkDirty() {
	l.dirty = true

	if l.parent != nil {
		l.parent.MarkDirty()
	}
}

func (l *Limit) Destroy() {
	if l.drawableObject.texture != nil {
		l.drawableObject.texture.Destroy()
	}

	l.limitParams.Child.Destroy()
}
