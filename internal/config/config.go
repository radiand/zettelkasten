package config

import "errors"
import "os"

import "github.com/BurntSushi/toml"

// Config represents global, application wide options.
type Config struct {
	ZettelkastenDir  string `toml:"zettelkasten_dir"`
	DefaultWorkspace string `toml:"default_workspace"`
}

// NewConfig creates config with default values.
func NewConfig() Config {
       return Config{
               ZettelkastenDir:  "~/vault/zettelkasten",
               DefaultWorkspace: "main",
       }
}

// GetConfig unmarshalls Config from array of bytes.
func GetConfig(content []byte) (Config, error) {
	var config Config
	err := toml.Unmarshal(content, &config)
	if err != nil {
		return Config{}, errors.Join(err, errors.New("Cannot get config"))
	}
	return config, nil
}

// GetConfigFromFile unmarshalls Config from file.
func GetConfigFromFile(path string) (Config, error) {
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return Config{}, errors.Join(readErr, errors.New("Cannot get config"))
	}
	config, marshalErr := GetConfig(content)
	if marshalErr != nil {
		return Config{}, errors.Join(marshalErr, errors.New("Cannot get config"))
	}
	return config, nil
}

// PutConfigToFile marshalls Config to file.
func PutConfigToFile(path string, config Config) error {
	content, err := toml.Marshal(config)
	if err != nil {
		return errors.Join(err, errors.New("Cannot marshall config"))
	}

	err = os.WriteFile(path, content, 0644)
	if err != nil {
		return errors.Join(err, errors.New("Cannot save config to file"))
	}
	return nil
}
