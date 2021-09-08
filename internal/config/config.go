package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func GetConfig(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)

	var config Config

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()

	return &config
}

type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}
