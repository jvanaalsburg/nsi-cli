package config

import (
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Auth struct {
		Email string `toml:"email"`
		Token string `toml:"token"`
	}
}

func (c Config) Save() error {
	data, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.Create("/etc/nsi/nsi-cli.toml")
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(data)

	return nil
}
