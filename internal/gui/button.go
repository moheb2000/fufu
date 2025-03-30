package gui

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Button struct {
	parent          Widget
	dirty           bool
	renderer        *sdl.Renderer
	buttonParams    *ButtonParams
	drawableObject  *DrawableObject
	backgroundColor sdl.Color
	onClick         func()
}

type ButtonParams struct {
	Value                *Text
	Color                sdl.Color
	ColorHover           sdl.Color
	BackgroundColor      sdl.Color
	BackgroundColorHover sdl.Color
	Width                int32
	Height               int32
}

func NewButton(renderer *sdl.Renderer, p *ButtonParams) (*Button, error) {
	b := Button{
		renderer:        renderer,
		dirty:           true,
		buttonParams:    p,
		backgroundColor: p.BackgroundColor,
		drawableObject:  &DrawableObject{},
	}

	// Change value text color to provided value in button params
	b.buttonParams.Value.setColor(b.buttonParams.Color)

	// Run setLimit to wrap text inside the button
	b.setLimit(int(b.buttonParams.Width))

	err := b.updateTexture()
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func (b *Button) updateTexture() error {
	// Check if there is a need to update the texture or not
	if !b.dirty {
		return nil
	}

	// If draw object is not nil, destroy the old texture first
	if b.drawableObject.texture != nil {
		b.drawableObject.texture.Destroy()
	}

	// Create a texture two combine character texture and value texture in one.
	texture, err := b.renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET, b.buttonParams.Width, b.buttonParams.Height)
	if err != nil {
		return err
	}

	b.renderer.SetRenderTarget(texture)

	// Making background transparent without lowering text quality
	texture.SetBlendMode(BLENDMOD_ONE)
	b.renderer.SetDrawColor(0, 0, 0, 0)
	b.renderer.Clear()

	vdo, err := b.buttonParams.Value.Draw()
	if err != nil {
		return err
	}

	// Set background color of the button
	b.renderer.SetDrawColor(b.backgroundColor.R, b.backgroundColor.G, b.backgroundColor.B, b.backgroundColor.A)
	b.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: b.buttonParams.Width, H: b.buttonParams.Height})

	// Make text inside the button center
	b.renderer.Copy(vdo.texture, nil, &sdl.Rect{X: b.buttonParams.Width/2 - vdo.W/2, Y: b.buttonParams.Height/2 - vdo.H/2, W: vdo.W, H: vdo.H})

	// Set render target back to nil
	b.renderer.SetRenderTarget(nil)

	b.drawableObject.texture = texture
	b.drawableObject.W = b.buttonParams.Width
	b.drawableObject.H = b.buttonParams.Height

	b.dirty = false

	return nil
}

func (b *Button) Draw() (*DrawableObject, error) {
	err := b.updateTexture()
	if err != nil {
		return nil, err
	}

	return b.drawableObject, nil
}

// HandleEvent handles mouse click button and the color hover change
func (b *Button) HandleEvent(event sdl.Event) {
	// Check the mouse position is inside the button
	if b.isMouseInside() {
		// Handle left click button
		switch e := event.(type) {
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONUP && e.Button == sdl.BUTTON_LEFT {
				if b.onClick != nil {
					b.onClick()
				}
			}
		}

		// When mouse is inside the button, check if b.backgroundColor is not equal to b.buttonParams.BackgroundColorHover and if they are not equal make them equal and mark button dirty
		if b.backgroundColor != b.buttonParams.BackgroundColorHover {
			b.backgroundColor = b.buttonParams.BackgroundColorHover
			b.buttonParams.Value.setColor(b.buttonParams.ColorHover)
			b.MarkDirty()
		}
	} else {
		// When mouse is outside the button, check if b.backgroundColor is not equal to b.buttonParams.BackgroundColor and if they are not equal make them equal and mark button dirty
		if b.backgroundColor != b.buttonParams.BackgroundColor {
			b.backgroundColor = b.buttonParams.BackgroundColor
			b.buttonParams.Value.setColor(b.buttonParams.Color)
			b.MarkDirty()
		}
	}

	b.buttonParams.Value.HandleEvent(event)
}

// OnClick gets a function as parameter and set it as onclick fallback
func (b *Button) OnClick(fn func()) {
	b.onClick = fn
}

// makeParent change the parent field to the provided argument
func (b *Button) makeParent(parent Widget) {
	b.parent = parent
}

// getParent returns the button widget parent
func (b *Button) getParent() Widget {
	return b.parent
}

// setLimit checks if the width of the button is bigger than the limit and if so, change the width to the new limit value and then run setLimit for child widget with the new width value
func (b *Button) setLimit(limit int) {
	// TODO: Add padding
	if limit != 0 && b.buttonParams.Width > int32(limit) {
		b.buttonParams.Width = int32(limit)
	}

	b.buttonParams.Value.setLimit(int(b.buttonParams.Width))

}

// MarkDirty changes the dirty parameter of widget to true. This function needs to be called if user wants to update the widget
func (b *Button) MarkDirty() {
	b.dirty = true

	if b.parent != nil {
		b.parent.MarkDirty()
	}
}

// isMouseInside checks if the mouse position is inside the button or not and returns a boolean value
func (b *Button) isMouseInside() bool {
	// Get mouse position
	mouseX, mouseY, _ := sdl.GetMouseState()
	var buttonX, buttonY int32
	// Get button position with adding all x and y params of drawable object for button parents
	for parent := b.parent; parent != nil; parent = parent.getParent() {
		do, _ := parent.Draw()

		buttonX += do.x
		buttonY += do.y

		// TODO: change the inside map for scrollable area also
		if l, ok := parent.(*List); ok {
			var previousHeight int32
			for _, widget := range l.listParams.Children {
				do, _ := widget.Draw()

				if b == widget {
					buttonY += previousHeight
					break
				}

				previousHeight += do.H + do.y + l.listParams.Spacing
			}
		}
	}

	// Check if the mouse position is inside the button
	if mouseX >= buttonX && mouseX <= buttonX+b.drawableObject.W && mouseY >= buttonY && mouseY <= buttonY+b.drawableObject.H {
		return true
	}

	return false
}

// Destroy cleans up the memory
func (b *Button) Destroy() {
	// Destroy drawable object
	if b.drawableObject.texture != nil {
		b.drawableObject.texture.Destroy()
	}

	// Destroy widget children
	b.buttonParams.Value.Destroy()
}
