package main

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Background struct {
	texture *sdl.Texture
}

func newBackground(renderer *sdl.Renderer, path string) (*Background, error) {
	bg := Background{}
	bgTexture, err := img.LoadTexture(renderer, path)
	if err != nil {
		return nil, err
	}
	bg.texture = bgTexture

	return &bg, nil
}

func (bg *Background) draw(renderer *sdl.Renderer) error {
	return renderer.Copy(bg.texture, nil, nil)
}

func (bg *Background) Destroy() {
	bg.texture.Destroy()
}
