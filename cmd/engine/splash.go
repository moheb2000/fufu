package main

import (
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Splash struct {
	renderer     *sdl.Renderer
	splashParams *SplashParams
	texture      *sdl.Texture
	opacity      float64
	timer        time.Time
}

type SplashParams struct {
	Path     string
	Color    sdl.Color
	Duration time.Duration
	App      *Application
}

func newSplash(renderer *sdl.Renderer, p *SplashParams) (*Splash, error) {
	s := Splash{
		renderer:     renderer,
		splashParams: p,
		opacity:      1,
	}

	splashTexture, err := img.LoadTexture(renderer, s.splashParams.Path)
	if err != nil {
		return nil, err
	}
	s.texture = splashTexture

	s.splashParams.App.am.Add(s.FadeIn())

	return &s, nil
}

func (s *Splash) draw() (bool, error) {
	resolution, err := s.splashParams.App.getResolution()
	if err != nil {
		return true, err
	}

	if s.opacity == 1 && s.timer.IsZero() {
		s.timer = time.Now()
	}

	if s.opacity == 1 && time.Since(s.timer) > s.splashParams.Duration {
		s.splashParams.App.am.Add(s.FadeOut())
	}

	s.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	s.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
	s.texture.SetAlphaMod(uint8(s.opacity * 255))

	s.renderer.SetDrawColor(s.splashParams.Color.R, s.splashParams.Color.G, s.splashParams.Color.B, uint8(s.opacity*255))
	s.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(resolution.X), H: int32(resolution.Y)})

	_, _, w, h, _ := s.texture.Query()
	s.renderer.Copy(s.texture, nil, &sdl.Rect{X: int32Abs(int32(resolution.X)/2 - w/2), Y: int32Abs(int32(resolution.Y)/2 - h/2), W: w, H: h})

	return s.opacity == 0, nil
}

func (s *Splash) FadeIn() func(time.Duration) bool {
	s.opacity = 0

	return func(dt time.Duration) bool {
		if s.opacity < 1 {
			s.opacity += dt.Seconds()
		}

		// It will ensure opacity not exceed from 1
		if s.opacity > 1 {
			s.opacity = 1
		}

		return s.opacity == 1
	}
}

func (s *Splash) FadeOut() func(time.Duration) bool {
	s.opacity = 1

	return func(dt time.Duration) bool {
		if s.opacity > 0 {
			s.opacity -= dt.Seconds()
		}

		// It will ensure opacity not exceed from 1
		if s.opacity < 0 {
			s.opacity = 0
		}

		return s.opacity == 0
	}
}

func (s *Splash) Destroy() {
	if s.texture != nil {
		s.texture.Destroy()
	}
}
