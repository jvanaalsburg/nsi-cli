package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml"
)

var CONFIG_FILE = "/etc/nsi/nsi-cli.toml"

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

	f, err := os.Create(CONFIG_FILE)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(data)

	return nil
}

func LoadConfig() (Config, error) {
	f, err := os.Open(CONFIG_FILE)
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
