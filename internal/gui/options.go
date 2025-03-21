package gui

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

type Options struct {
	parent         Widget
	dirty          bool
	done           bool
	renderer       *sdl.Renderer
	optionsParams  *OptionsParams
	drawableObject *DrawableObject
}

type OptionsParams struct {
	Options *List
	Result  *int
}

func NewOptions(renderer *sdl.Renderer, p *OptionsParams) (*Options, error) {
	o := Options{
		renderer:       renderer,
		optionsParams:  p,
		dirty:          true,
		done:           false,
		drawableObject: &DrawableObject{},
	}

	p.Options.makeParent(&o)

	err := o.updateTexture()
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (o *Options) updateTexture() error {
	// Check if there is a need to update the texture or not
	if !o.dirty {
		return nil
	}

	// If draw object is not nil, destroy the old texture first
	if o.drawableObject.texture != nil {
		o.drawableObject.texture.Destroy()
	}

	o.drawableObject, _ = o.optionsParams.Options.Draw()

	o.dirty = false

	return nil
}

func (o *Options) Draw() (*DrawableObject, error) {
	err := o.updateTexture()
	if err != nil {
		return nil, err
	}

	return o.drawableObject, nil
}

func (o *Options) HandleEvent(event sdl.Event) {
	// Handle user input here
	if o.done {
		return
	}

	switch e := event.(type) {
	case *sdl.KeyboardEvent:
		if e.Type == sdl.KEYUP {
			for i := range o.optionsParams.Options.listParams.Children {
				if e.Keysym.Sym == sdl.GetKeyFromName(strconv.Itoa(i+1)) {
					if *o.optionsParams.Result == 0 {
						*o.optionsParams.Result = i + 1
					}
					o.done = true

					break
				}
			}
		}
	}

	o.optionsParams.Options.HandleEvent(event)
}

func (o *Options) makeParent(parent Widget) {
	o.parent = parent
}

func (o *Options) getParent() Widget {
	return o.parent
}

func (o *Options) setLimit(limit int) {
	o.optionsParams.Options.setLimit(limit)
}

func (o *Options) MarkDirty() {
	o.dirty = true

	if o.parent != nil {
		o.parent.MarkDirty()
	}
}

func (o *Options) Destroy() {
	// Destroy drawable object
	if o.drawableObject.texture != nil {
		o.drawableObject.texture.Destroy()
	}

	o.optionsParams.Options.Destroy()
}
