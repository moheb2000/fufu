package gui

import (
	"github.com/veandco/go-sdl2/ttf"
)

type FontManager struct {
	fonts map[string]map[int]*ttf.Font
}

func NewFontManager() *FontManager {
	return &FontManager{fonts: make(map[string]map[int]*ttf.Font)}
}

func (fm *FontManager) LoadFont(name string, path string, size int) (*ttf.Font, error) {
	if _, exists := fm.fonts[name]; !exists {
		fm.fonts[name] = make(map[int]*ttf.Font)
	}

	if font, exists := fm.fonts[name][size]; exists {
		return font, nil
	}

	font, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, err
	}

	fm.fonts[name][size] = font
	return font, nil
}

func (fm *FontManager) GetFont(name string, size int) *ttf.Font {
	if sizes, exists := fm.fonts[name]; exists {
		if font, exists := sizes[size]; exists {
			return font
		}
	}

	return nil
}

func (fm *FontManager) Close() {
	for _, sizes := range fm.fonts {
		for _, font := range sizes {
			font.Close()
		}
	}
}
