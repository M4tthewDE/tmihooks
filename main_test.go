package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/structs"
	"github.com/stretchr/testify/assert"
)

func TestApplication(t *testing.T) {
	t.Parallel()

	testServer := TestServer{
		t: t,
	}

	go testServer.startTestClient()

	testServer.registerWebhook()

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
	http.HandleFunc("/chat", ts.chat)

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

	_, err = w.Write([]byte(confirmation.Challenge))
	if err != nil {
		panic(err)
	}
}

func (ts *TestServer) chat(w http.ResponseWriter, req *http.Request) {
	var msg twitch.PrivateMessage

	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		panic(err)
	}

	log.Println(msg.Channel, msg.Message)
}

func (ts *TestServer) registerWebhook() {
	config := config.GetConfig("test_config.yml")
	webhook := structs.Webhook{
		Channels:    []string{"tmiloadtesting2", "twitchmedia_qs_10", "nmplol", "gtawiseguy", "quin69"},
		URI:         "http://localhost:7070/chat",
		RegisterURI: "http://localhost:7070/register",
		Nonce:       "penis123",
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
