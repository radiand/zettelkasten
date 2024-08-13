package main

import "os"
import "fmt"
import "github.com/BurntSushi/toml"

// Config represents global, application wide options.
type Config struct {
	// Path is a root directory where notes are stored.
	Path string `toml:"path"`
}

// GetConfig unmarshalls Config from array of bytes.
func GetConfig(content []byte) (Config, error) {
	var config Config
	err := toml.Unmarshal(content, &config)
	if err != nil {
		return Config{}, fmt.Errorf("Cannot get config due to: %w", err)
	}
	return config, nil
}

// GetConfigFromFile unmarshalls Config from file.
func GetConfigFromFile(path string) (Config, error) {
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return Config{}, fmt.Errorf("Cannot get config due to: %w", readErr)
	}
	config, marshalErr := GetConfig(content)
	if marshalErr != nil {
		return Config{}, fmt.Errorf("Cannot get config due to: %w", marshalErr)
	}
	return config, nil
}
