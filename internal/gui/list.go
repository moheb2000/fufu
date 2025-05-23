package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type List struct {
	parent         Widget
	dirty          bool
	renderer       *sdl.Renderer
	listParams     *ListParams
	drawableObject *DrawableObject
}

type ListParams struct {
	Children []Widget
	Spacing  int32
}

func NewList(renderer *sdl.Renderer, p *ListParams) (*List, error) {
	l := List{
		renderer:       renderer,
		listParams:     p,
		dirty:          true,
		drawableObject: &DrawableObject{},
	}

	for _, widget := range l.listParams.Children {
		widget.makeParent(&l)
	}

	err := l.updateTexture()
	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (l *List) updateTexture() error {
	if !l.dirty {
		return nil
	}

	if l.drawableObject.texture != nil {
		l.drawableObject.texture.Destroy()
	}

	var lw int32
	var lh int32
	for _, widget := range l.listParams.Children {
		do, _ := widget.Draw()
		if lw < do.W {
			lw = do.W + do.x
		}

		lh += do.H + do.y + l.listParams.Spacing
	}
	lh -= l.listParams.Spacing

	// Fix error: Texture dimentions can't be zero
	if lw <= 0 {
		lw = 1
	}
	if lh <= 0 {
		lh = 1
	}

	texture, err := l.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, lw, lh)
	if err != nil {
		return err
	}

	l.renderer.SetRenderTarget(texture)

	// Making background transparent without reducing quality
	texture.SetBlendMode(BLENDMOD_ONE)
	l.renderer.SetDrawColor(0, 0, 0, 0)
	l.renderer.Clear()

	var previousHeight int32
	for _, widget := range l.listParams.Children {
		do, _ := widget.Draw()

		l.renderer.Copy(do.texture, nil, &sdl.Rect{X: do.x, Y: do.y + previousHeight, W: do.W, H: do.H})
		previousHeight += do.H + do.y + l.listParams.Spacing
	}

	l.renderer.SetRenderTarget(nil)

	l.drawableObject.W = lw
	l.drawableObject.H = lh
	l.drawableObject.texture = texture

	// Reset dirty flag
	l.dirty = false

	return nil
}

func (l *List) Draw() (*DrawableObject, error) {
	err := l.updateTexture()
	if err != nil {
		return nil, err
	}

	return l.drawableObject, nil
}

func (l *List) HandleEvent(event sdl.Event) {
	for _, widget := range l.listParams.Children {
		widget.HandleEvent(event)
	}
}

func (l *List) makeParent(parent Widget) {
	l.parent = parent
}

func (l *List) getParent() Widget {
	return l.parent
}

func (l *List) setLimit(limit int) {
	for _, widget := range l.listParams.Children {
		widget.setLimit(limit)
	}
}

func (l *List) MarkDirty() {
	l.dirty = true

	if l.parent != nil {
		l.parent.MarkDirty()
	}
}

func (l *List) AddWidget(w Widget) {
	l.listParams.Children = append([]Widget{w}, l.listParams.Children...)

	w.makeParent(l)
	l.MarkDirty()
}

func (l *List) RemoveLastWidget() {
	l.listParams.Children[0].Destroy()
	l.listParams.Children = l.listParams.Children[1:]
	l.MarkDirty()
}

func (l *List) Destroy() {
	if l.drawableObject.texture != nil {
		l.drawableObject.texture.Destroy()
	}

	for _, widget := range l.listParams.Children {
		widget.Destroy()
	}
}
