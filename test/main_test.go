package test

import (
	"testing"

	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/db"
	"github.com/m4tthewde/tmihooks/internal/tmi"
)

func RunMainApplication() {
	config := config.GetConfig("../config.yml")

	dbHandler := db.DatabaseHandler{
		Config: config,
	}

	// clear mongodb and redis
	db.Clear()
	err := dbHandler.Clear()
	if err != nil {
		panic(err)
	}

	reader := tmi.NewReader(config)

	go reader.Read()

	server := api.NewServer(config, reader)
	go server.Run()
}

func TestRegister(t *testing.T) {
	RunMainApplication()

	testServer := NewTestServer(t, REGISTER)

	go testServer.StartTestClient()

	testServer.RegisterWebhook()
	<-testServer.StopChan
}
