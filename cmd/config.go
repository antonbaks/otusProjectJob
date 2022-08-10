package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	CacheSize int
	UploadDir string
	HTTPHost  string
	HTTPPort  string
	MaxWidth  uint
	MinWidth  uint
	MaxHeight uint
	MinHeight uint
}

func NewConfig(path string) *Config {
	var c Config

	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}
