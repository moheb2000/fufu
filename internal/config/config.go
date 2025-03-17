// config file contains logic to open and access to engine configurations
package config

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct is a model for data in config.json
type Config struct {
	Title       string
	FPS         int
	FullScreen  bool
	Resolution  int
	DefaultFont string
	DialogPanel struct {
		Direction string
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
		Title:       "Fufu Visual Novel Engine",
		FPS:         30,
		FullScreen:  false,
		Resolution:  1080,
		DefaultFont: "assets/UbuntuSans-Regular.ttf",
		DialogPanel: struct{ Direction string }{
			Direction: "left",
		},
	}

	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	if cfg.DialogPanel.Direction != "left" && cfg.DialogPanel.Direction != "right" {
		log.Println("Dialog panel direction is invalid. Engine use \"left\" as fallback direction")
		cfg.DialogPanel.Direction = "left"
	}

	return &cfg, nil
}
