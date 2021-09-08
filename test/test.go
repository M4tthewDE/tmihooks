package main

import (
	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
)

func TestApplication() {
	config := config.GetConfig("test_config.yml")

	server := api.NewServer(config)
	server.Run()
}
