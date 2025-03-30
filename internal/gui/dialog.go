package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Dialog struct {
	parent         Widget
	dirty          bool
	renderer       *sdl.Renderer
	dialogParams   *DialogParams
	drawableObject *DrawableObject
}

type DialogParams struct {
	Character *Text
	Value     *Text
}

// NewDialog returns a new Dialog struct widget with the provided parameters
func NewDialog(renderer *sdl.Renderer, p *DialogParams) (*Dialog, error) {
	// Create a base Dialog struct. Setting dirty to true will cause the texture to be created for the first time
	d := Dialog{
		renderer:       renderer,
		dialogParams:   p,
		dirty:          true,
		drawableObject: &DrawableObject{},
	}

	// Make dialog parent for children widgets
	p.Character.makeParent(&d)
	p.Value.makeParent(&d)

	// TODO: Check if there is a need for updating texture in initialize functions or not. Because instead it will tun in the first run of main loop becuase of dirty is true
	err := d.updateTexture()
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// updateTexture updates the texture of the widget in every frame if dirty is true
func (d *Dialog) updateTexture() error {
	// Check if there is a need to update the texture or not
	if !d.dirty {
		return nil
	}

	// If draw object is not nil, destroy the old texture first
	if d.drawableObject.texture != nil {
		d.drawableObject.texture.Destroy()
	}

	// character and value text widgets
	ctdo, _ := d.dialogParams.Character.Draw()
	vtdo, _ := d.dialogParams.Value.Draw()

	// Create a texture two combine character texture and value texture in one.
	texture, err := d.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, ctdo.W+vtdo.W, vtdo.H)
	if err != nil {
		return err
	}

	d.renderer.SetRenderTarget(texture)

	// Making background transparent without lowering text quality
	texture.SetBlendMode(BLENDMOD_ONE)
	d.renderer.SetDrawColor(0, 0, 0, 0)
	d.renderer.Clear()

	// Add character texture to the combined texture
	d.renderer.Copy(ctdo.texture, nil, &sdl.Rect{X: 0, Y: 0, W: ctdo.W, H: ctdo.H})

	// Add value texture to the combined texture
	d.renderer.Copy(vtdo.texture, nil, &sdl.Rect{X: ctdo.W, Y: 0, W: vtdo.W, H: vtdo.H})

	// Set render target back to nil
	d.renderer.SetRenderTarget(nil)

	// Update w, h and texture of drawable object
	d.drawableObject.W = ctdo.W + vtdo.W
	d.drawableObject.H = vtdo.H
	d.drawableObject.texture = texture

	// Reset dirty flag
	d.dirty = false

	return nil
}

// Draw returns a drawable object. If texture must rerender, it will happen here
func (d *Dialog) Draw() (*DrawableObject, error) {
	err := d.updateTexture()
	if err != nil {
		return nil, err
	}

	return d.drawableObject, nil
}

func (d *Dialog) HandleEvent(event sdl.Event) {
	d.dialogParams.Character.HandleEvent(event)
	d.dialogParams.Value.HandleEvent(event)
}

func (d *Dialog) makeParent(parent Widget) {
	d.parent = parent
}

func (d *Dialog) getParent() Widget {
	return d.parent
}

func (d *Dialog) setLimit(limit int) {
	d.dialogParams.Value.setLimit(limit - int(d.dialogParams.Character.drawableObject.W))
}

// MarkDirty changes the dirty parameter of widget to true. This function needs to be called if user wants to update the widget
func (d *Dialog) MarkDirty() {
	d.dirty = true

	if d.parent != nil {
		d.parent.MarkDirty()
	}
}

// Destroy cleans up the textures from the memory
func (d *Dialog) Destroy() {
	// Destroy drawable object
	if d.drawableObject.texture != nil {
		d.drawableObject.texture.Destroy()
	}

	// Destroy widget  children
	d.dialogParams.Character.Destroy()
	d.dialogParams.Value.Destroy()
}
