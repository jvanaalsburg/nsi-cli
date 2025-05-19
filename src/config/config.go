package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

func configFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(configDir, "nsi", "nsi-cli.toml")
	_, err = os.Stat(path)
	if err != nil {
		return "", err
	}

	return path, nil
}

type Config struct {
	Api  ApiConfig  `toml:"api"`
	Auth AuthConfig `toml:"auth"`
}

type ApiConfig struct {
	UrlRoot string `toml:"url_root"`
}

type AuthConfig struct {
	Email string `toml:"email"`
	Token string `toml:"token"`
}

func (c Config) Save() error {
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	filename, err := configFile()
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(data)

	return nil
}

func LoadConfig() (Config, error) {
	filename, err := configFile()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
