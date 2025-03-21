package gui

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Text struct {
	parent         Widget
	dirty          bool
	renderer       *sdl.Renderer
	textParams     *TextParams
	drawableObject *DrawableObject
	opacity        float64
}

type TextParams struct {
	Value string
	Color sdl.Color
	Font  *ttf.Font
	limit int
}

// Create a new Text widget and returns it
func NewText(renderer *sdl.Renderer, p *TextParams) (*Text, error) {
	// Making dirty to true force widget to be rendered for the first time
	t := Text{
		renderer:       renderer,
		textParams:     p,
		dirty:          true,
		drawableObject: &DrawableObject{},
		opacity:        1,
	}

	err := t.updateTexture()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Draw returns a drawable object. If the widget needs to update, this will happen here
func (t *Text) Draw() (*DrawableObject, error) {
	err := t.updateTexture()
	if err != nil {
		return nil, err
	}

	return t.drawableObject, nil
}

func (t *Text) HandleEvent(event sdl.Event) {}

func (t *Text) FadeIn() func(time.Duration) bool {
	t.opacity = 0

	return func(dt time.Duration) bool {
		if t.opacity < 1 {
			t.opacity += dt.Seconds()
			t.MarkDirty()
		}

		if t.opacity > 1 {
			t.opacity = 1
		}

		return t.opacity == 1
	}
}

// updateTexture updates the texture of the widget in every frame if dirty is true
func (t *Text) updateTexture() error {
	// Check if texture needs to update
	if !t.dirty {
		return nil
	}

	// Destroy the old texture to prevent leaking memory
	if t.drawableObject.texture != nil {
		t.drawableObject.texture.Destroy()
	}

	t.textParams.Color.A = uint8(t.opacity * 255)

	// Make a text texture with font and specified color
	surface, err := t.textParams.Font.RenderUTF8BlendedWrapped(t.textParams.Value, t.textParams.Color, t.textParams.limit)
	if err != nil {
		return err
	}
	defer surface.Free()

	t.drawableObject.w = surface.W
	t.drawableObject.h = surface.H

	texture, err := t.renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return err
	}

	t.drawableObject.texture = texture

	// Reset dirty flag
	t.dirty = false

	return nil
}

func (t *Text) makeParent(parent Widget) {
	t.parent = parent
}

func (t *Text) getParent() Widget {
	return t.parent
}

func (t *Text) setLimit(limit int) {
	t.textParams.limit = limit
	t.MarkDirty()
}

// MarkDirty changes the dirty parameter of widget to true. This function needs to be called if user wants to update the widget
func (t *Text) MarkDirty() {
	t.dirty = true

	if t.parent != nil {
		// fmt.Println(t.textParams.limit)
		t.parent.MarkDirty()
	}
}

// Cleanup the memory
func (t *Text) Destroy() {
	if t.drawableObject.texture != nil {
		t.drawableObject.texture.Destroy()
	}
}
