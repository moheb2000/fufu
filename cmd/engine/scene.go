package main

import (
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Background struct {
	backgroundParams *BackgroundParams
	texture          *sdl.Texture
	renderer         *sdl.Renderer
	opacity          float64
}

type BackgroundParams struct {
	Origin *Origin
	Path   string
	App    *Application
}

type Origin struct {
	X string
	Y string
}

func newBackground(renderer *sdl.Renderer, p *BackgroundParams) (*Background, error) {
	bg := Background{
		backgroundParams: p,
		renderer:         renderer,
		opacity:          1,
	}
	bgTexture, err := img.LoadTexture(renderer, bg.backgroundParams.Path)
	if err != nil {
		return nil, err
	}
	bg.texture = bgTexture

	return &bg, nil
}

func (bg *Background) draw() error {
	resolution, err := bg.backgroundParams.App.getResolution()
	if err != nil {
		return err
	}

	// This is equivalent to bgTextRect.W
	panelWidth := int32(float64(resolution.X) * bg.backgroundParams.App.cfg.DialogPanel.Width)

	_, _, tw, th, _ := bg.texture.Query()
	dw := int32(resolution.X) - panelWidth
	dh := int32(resolution.Y)

	var sw, sh int32
	if tw/th >= dw/dh {
		sw = th * dw / dh
		sh = th
	} else {
		sw = tw
		sh = tw * dh / dw
	}

	var x, y int32
	if bg.backgroundParams.App.cfg.DialogPanel.Direction == "left" {
		x = panelWidth
	}

	var xOffset, yOffset int32
	switch bg.backgroundParams.Origin.X {
	case "left":
		xOffset = 0
	case "center":
		xOffset = int32Abs(sw/2 - dw/2)
	case "right":
		xOffset = int32Abs(sw - dw)
	}

	switch bg.backgroundParams.Origin.Y {
	case "top":
		yOffset = 0
	case "center":
		xOffset = int32Abs(sh/2 - dh/2)
	case "bottom":
		xOffset = int32Abs(sh - dh)
	}

	bg.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	bg.texture.SetAlphaMod(uint8(bg.opacity * 255))

	return bg.renderer.Copy(bg.texture, &sdl.Rect{X: xOffset, Y: yOffset, W: sw, H: sh}, &sdl.Rect{X: x, Y: y, W: dw, H: dh})
}

func (bg *Background) FadeIn() func(time.Duration) bool {
	bg.opacity = 0

	return func(dt time.Duration) bool {
		if bg.opacity < 1 {
			bg.opacity += dt.Seconds()
		}

		// It will ensure opacity not exceed from 1
		if bg.opacity > 1 {
			bg.opacity = 1
		}

		return bg.opacity == 1
	}
}

func (bg *Background) Destroy() {
	if bg.texture != nil {
		bg.texture.Destroy()
	}
}
