package main

import (
	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/db"
	"github.com/m4tthewde/tmihooks/internal/tmi"
)

func main() {
	// clear redis
	db.Clear()

	config := config.GetConfig("config.yml")

	reader := tmi.NewReader(config)

	go reader.Read()

	server := api.NewServer(config, reader)
	server.Run()
}
