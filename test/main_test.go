package test

import (
	"testing"

	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/db"
	"github.com/m4tthewde/tmihooks/internal/tmi"
)

func RunMainApplication() {
	// run main
	// clear redis
	db.Clear()

	config := config.GetConfig("../config.yml")

	dbHandler := db.DatabaseHandler{
		Config: config,
	}

	dbHandler.Clear()

	reader := tmi.NewReader(config)

	go reader.Read()

	server := api.NewServer(config, reader)
	go server.Run()
}

func TestRegister(t *testing.T) {
	t.Parallel()
	RunMainApplication()

	testServer := NewTestServer(t)

	go testServer.StartTestClient()

	testServer.RegisterWebhook()

	for {
		select {}
	}
}
