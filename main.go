package main

import (
	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/tmi"
)

func main() {
	config := config.GetConfig("config.yml")

	reader := tmi.NewReader()

	go reader.Read()

	server := api.NewServer(config, reader)
	server.Run()
}
