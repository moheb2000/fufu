package gui

import (
	"github.com/veandco/go-sdl2/ttf"
)

type FontManager struct {
	fonts map[string]map[int]*ttf.Font
}

// NewFontManager returns a new FontManager pointer and make an empty map for fonts field
func NewFontManager() *FontManager {
	return &FontManager{fonts: make(map[string]map[int]*ttf.Font)}
}

// LoadFont gets font path and size, check if it is exists or not and if it's not, add it to fonts map as a cache storage
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

// GetFont returns a font based on name and size and if it doesn't exist, returns nil
func (fm *FontManager) GetFont(name string, size int) *ttf.Font {
	if sizes, exists := fm.fonts[name]; exists {
		if font, exists := sizes[size]; exists {
			return font
		}
	}

	return nil
}

// Free memory from all cached fonts
func (fm *FontManager) Close() {
	for _, sizes := range fm.fonts {
		for _, font := range sizes {
			font.Close()
		}
	}
}
