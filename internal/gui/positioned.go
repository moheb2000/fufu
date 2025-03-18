package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Positioned struct {
	parent           Widget
	dirty            bool
	renderer         *sdl.Renderer
	positionedParams *PositionedParams
	drawableObject   *DrawableObject
}

type PositionedParams struct {
	X     int32
	Y     int32
	Child Widget
}

func NewPositioned(renderer *sdl.Renderer, p *PositionedParams) (*Positioned, error) {
	pos := Positioned{
		renderer:         renderer,
		positionedParams: p,
		dirty:            true,
		drawableObject:   &DrawableObject{},
	}

	p.Child.makeParent(&pos)

	err := pos.updateTexture()
	if err != nil {
		return nil, err
	}

	return &pos, nil
}

func (p *Positioned) updateTexture() error {
	if !p.dirty {
		return nil
	}

	do, err := p.positionedParams.Child.Draw()
	if err != nil {
		return err
	}

	p.drawableObject = do
	p.drawableObject.x = p.positionedParams.X
	p.drawableObject.y = p.positionedParams.Y

	p.dirty = false

	return nil
}

func (p *Positioned) Draw() (*DrawableObject, error) {
	err := p.updateTexture()
	if err != nil {
		return nil, err
	}

	return p.drawableObject, nil
}

func (p *Positioned) HandleEvent(event sdl.Event) {
	p.positionedParams.Child.HandleEvent(event)
}

func (p *Positioned) makeParent(parent Widget) {
	p.parent = parent
}

func (p *Positioned) getParent() Widget {
	return p.parent
}

func (p *Positioned) setLimit(limit int) {
	p.positionedParams.Child.setLimit(limit)
}

func (p *Positioned) MarkDirty() {
	p.dirty = true

	if p.parent != nil {
		p.parent.MarkDirty()
	}
}

func (p *Positioned) Destroy() {
	if p.drawableObject.texture != nil {
		p.drawableObject.texture.Destroy()
	}

	p.positionedParams.Child.Destroy()
}
