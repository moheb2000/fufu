package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ScrollableArea struct {
	parent               Widget
	dirty                bool
	renderer             *sdl.Renderer
	scrollableAreaParams *ScrollableAreaParams
	drawableObject       *DrawableObject
	scroll               int32
}

type ScrollableAreaParams struct {
	H          int32
	Child      Widget
	ScrollStep int32
}

func NewScrollableArea(renderer *sdl.Renderer, p *ScrollableAreaParams) (*ScrollableArea, error) {
	s := ScrollableArea{
		renderer:             renderer,
		scrollableAreaParams: p,
		dirty:                true,
		drawableObject:       &DrawableObject{},
	}

	if s.scrollableAreaParams.ScrollStep <= 0 {
		s.scrollableAreaParams.ScrollStep = 5
	}

	p.Child.makeParent(&s)

	err := s.updateTexture()
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *ScrollableArea) updateTexture() error {
	if !s.dirty {
		return nil
	}

	if s.drawableObject.texture != nil {
		s.drawableObject.texture.Destroy()
	}

	do, _ := s.scrollableAreaParams.Child.Draw()

	texture, err := s.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, do.w, s.scrollableAreaParams.H)
	if err != nil {
		return err
	}

	s.renderer.SetRenderTarget(texture)

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	s.renderer.SetDrawColor(255, 255, 255, 0)
	s.renderer.Clear()

	if s.scroll < 0 {
		s.scroll = 0
	} else if s.scroll > do.h-s.scrollableAreaParams.H {
		s.scroll = do.h - s.scrollableAreaParams.H
	}

	dst := sdl.Rect{X: do.x, Y: do.y, W: do.w, H: s.scrollableAreaParams.H}
	if do.h < s.scrollableAreaParams.H {
		dst.H = do.h
	}

	s.renderer.Copy(do.texture, &sdl.Rect{X: 0, Y: s.scroll, W: do.w, H: s.scrollableAreaParams.H}, &dst)

	s.renderer.SetRenderTarget(nil)

	s.drawableObject.texture = texture
	s.drawableObject.w = do.w
	s.drawableObject.h = s.scrollableAreaParams.H

	s.dirty = false

	return nil
}

func (s *ScrollableArea) Draw() (*DrawableObject, error) {
	err := s.updateTexture()
	if err != nil {
		return nil, err
	}

	return s.drawableObject, nil
}

func (s *ScrollableArea) HandleEvent(event sdl.Event) {
	switch e := event.(type) {
	case *sdl.MouseWheelEvent:
		if s.isMouseInside() {
			s.scroll -= e.Y * s.scrollableAreaParams.ScrollStep
			// TODO: Call MarkDirty if scroll actually changes to something and not the beginning and end of the scrollable area
			s.MarkDirty()
		}
	}

	s.scrollableAreaParams.Child.HandleEvent(event)
}

func (s *ScrollableArea) makeParent(parent Widget) {
	s.parent = parent
}

func (s *ScrollableArea) getParent() Widget {
	return s.parent
}

func (s *ScrollableArea) setLimit(limit int) {
	s.scrollableAreaParams.Child.setLimit(limit)
}

func (s *ScrollableArea) MarkDirty() {
	s.dirty = true

	if s.parent != nil {
		s.parent.MarkDirty()
	}
}

func (s *ScrollableArea) isMouseInside() bool {
	mouseX, mouseY, _ := sdl.GetMouseState()
	var scrollableX, scrollableY int32
	for parent := s.parent; parent != nil; parent = parent.getParent() {
		do, _ := parent.Draw()
		scrollableX += do.x
		scrollableY += do.y
	}

	if mouseX >= scrollableX && mouseX <= scrollableX+s.drawableObject.w && mouseY >= scrollableY && mouseY <= scrollableY+s.drawableObject.h {
		return true
	}

	return false
}

func (s *ScrollableArea) Destroy() {
	if s.drawableObject.texture != nil {
		s.drawableObject.texture.Destroy()
	}

	s.scrollableAreaParams.Child.Destroy()
}
