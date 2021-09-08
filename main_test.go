package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/m4tthewde/tmihooks/internal/api"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/structs"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	t.Parallel()

	config := config.GetConfig("test_config.yml")

	testServer := TestServer{
		t: t,
	}

	go testServer.startTestClient()

	server := api.NewServer(config)
	go server.Run()

	testServer.registerWebhook(config)

	for {
		select {}
	}
}

type TestServer struct {
	t       *testing.T
	webhook *structs.Webhook
}

func (ts *TestServer) startTestClient() {
	http.HandleFunc("/register", ts.register)

	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		panic(err)
	}
}

func (ts *TestServer) register(w http.ResponseWriter, req *http.Request) {
	var confirmation structs.Confirmation

	err := json.NewDecoder(req.Body).Decode(&confirmation)
	if err != nil {
		panic(err)
	}

	log.Println("received confirmation")

	s, _ := json.MarshalIndent(confirmation, "", " ")

	log.Println(string(s))

	assert.Equal(ts.t, ts.webhook.Nonce, confirmation.Nonce)
}

func (ts *TestServer) registerWebhook(config *config.Config) {
	webhook := structs.Webhook{
		Channels: []string{"matthewde", "gopherobot"},
		URI:      "http://localhost:7070/register",
		Nonce:    "penis123",
	}
	s, _ := json.MarshalIndent(webhook, "", " ")

	log.Println("registering webhook...")
	log.Println(string(s))

	ts.webhook = &webhook

	webhookJSON, err := json.Marshal(webhook)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	address := "http://localhost:" + config.Server.Port + "/register"

	req, err := http.NewRequestWithContext(ctx, "POST", address, bytes.NewBuffer(webhookJSON))
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	assert.Equal(ts.t, http.StatusOK, resp.StatusCode)
}
