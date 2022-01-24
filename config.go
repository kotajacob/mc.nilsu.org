package main

import (
	"github.com/BurntSushi/toml"
)

type config struct {
	Address    string
	Template   string
	MCAddress  string
	ModList    string
	CarpetList string
}

func (c *config) Load(p string) error {
	if _, err := toml.DecodeFile(p, c); err != nil {
		return err
	}
	return nil
}
