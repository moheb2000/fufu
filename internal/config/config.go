// config file contains logic to open and access to engine configurations
package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct is a model for data in config.json
type Config struct {
	Title            string
	FPS              int
	GameVersion      string
	FullScreen       bool
	Resolution       int
	DefaultFont      string
	DefaultTextColor string
	DialogPanel      struct {
		Direction string
		Color     string
		Width     float64
	}
	MainMenu struct {
		Color                string
		ColorHover           string
		BackgroundColor      string
		BackgroundColorHover string
		Background           string
	}
}

// Get returns a Config struct that has configs from user or the default one
func Get() (*Config, error) {
	// Open the config.json file in the root of project
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Here is the default configs for engine
	cfg := Config{
		Title:            "Fufu Visual Novel Engine",
		FPS:              30,
		GameVersion:      "undefined",
		FullScreen:       false,
		Resolution:       1080,
		DefaultFont:      "assets/UbuntuSans-Regular.ttf",
		DefaultTextColor: "#ffffff",
		DialogPanel: struct {
			Direction string
			Color     string
			Width     float64
		}{
			Direction: "left",
			Color:     "#202020",
			Width:     0.3,
		},
		MainMenu: struct {
			Color                string
			ColorHover           string
			BackgroundColor      string
			BackgroundColorHover string
			Background           string
		}{
			Color:                "#ffffff",
			ColorHover:           "#000000",
			BackgroundColor:      "#045147",
			BackgroundColorHover: "#ffffff",
			Background:           "",
		},
	}

	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.DialogPanel.Width < 0.1 || cfg.DialogPanel.Width > 1 {
		cfg.DialogPanel.Width = 0.3
	}

	if cfg.DialogPanel.Direction != "left" && cfg.DialogPanel.Direction != "right" {
		log.Println("Dialog panel direction is invalid. Engine use \"left\" as fallback direction")
		cfg.DialogPanel.Direction = "left"
	}

	return &cfg, nil
}
