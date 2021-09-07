package main

import (
	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
)

func main() {
	config := config.GetConfig()

	server := api.NewServer(config)
	server.Run()
}
