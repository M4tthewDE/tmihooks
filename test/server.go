package test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/structs"
	"github.com/stretchr/testify/assert"
)

type TestServer struct {
	t        *testing.T
	server   *http.Server
	webhook  *structs.Webhook
	stopChan chan os.Signal
}

func NewTestServer(t *testing.T) *TestServer {
	ts := TestServer{
		t:        t,
		webhook:  nil,
		server:   &http.Server{Addr: ":7070"},
		stopChan: make(chan os.Signal, 1),
	}
	signal.Notify(ts.stopChan, os.Interrupt)

	return &ts
}

func (ts *TestServer) StartTestClient() {
	http.HandleFunc("/register", ts.Register)
	http.HandleFunc("/chat", ts.chat)

	go func() {
		err := ts.server.ListenAndServe()
		if err != nil {
			log.Println("server closed")
		}
	}()

	<-ts.stopChan
	log.Println("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ts.server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func (ts *TestServer) Register(w http.ResponseWriter, req *http.Request) {
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

	assert.Equal(ts.t, "tmiloadtesting2", msg.Channel)
	log.Println("received message, attempting graceful shutdown")

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		panic(err)
	}

	p.Signal(os.Interrupt)
}

func (ts *TestServer) RegisterWebhook() {
	config := config.GetConfig("../test_config.yml")
	webhook := structs.Webhook{
		Channels:    []string{"tmiloadtesting2"},
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
