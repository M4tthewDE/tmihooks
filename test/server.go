package test

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

type Type int

const (
	REGISTER Type = iota
	DELETE   Type = iota
	INFINITE Type = iota
)

type Server struct {
	t        *testing.T
	config   *config.Config
	server   *http.Server
	webhook  *structs.Webhook
	StopChan chan int
	testType Type
}

func NewTestServer(t *testing.T, testType Type) *Server {
	config := config.GetConfig("../test_config.yml")

	ts := Server{
		t:        t,
		config:   config,
		webhook:  nil,
		server:   &http.Server{Addr: ":7070"},
		StopChan: make(chan int, 1),
		testType: testType,
	}

	return &ts
}

func (ts *Server) StartTestClient() {
	http.HandleFunc("/register", ts.Register)
	http.HandleFunc("/chat", ts.chat)

	err := ts.server.ListenAndServe()
	if err != nil {
		log.Println("server closed")
	}
}

func (ts *Server) Register(w http.ResponseWriter, req *http.Request) {
	var confirmation structs.Confirmation

	err := json.NewDecoder(req.Body).Decode(&confirmation)
	if err != nil {
		panic(err)
	}

	log.Println("received confirmation")

	s, _ := json.MarshalIndent(confirmation, "", " ")

	log.Println(string(s))

	if ts.testType == REGISTER {
		assert.Equal(ts.t, ts.webhook.Nonce, confirmation.Nonce)
	}

	_, err = w.Write([]byte(confirmation.Challenge))
	if err != nil {
		panic(err)
	}
}

func (ts *Server) chat(w http.ResponseWriter, req *http.Request) {
	var msg twitch.PrivateMessage

	err := json.NewDecoder(req.Body).Decode(&msg)
	if err != nil {
		panic(err)
	}

	if ts.testType == REGISTER {
		assert.Equal(ts.t, "tmiloadtesting2", msg.Channel)
		log.Println("received message, attempting graceful shutdown")

		// shutdown main server
		ts.ShutdownMainServer()

		// shutdown test server
		ts.StopChan <- 0
	}
}

func (ts *Server) RegisterWebhook() {
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

	address := "http://localhost:" + ts.config.Server.Port + "/register"

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

func (ts *Server) ShutdownMainServer() {
	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	address := "http://localhost:" + ts.config.Server.Port + "/shutdown"

	req, err := http.NewRequestWithContext(ctx, "POST", address, nil)
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
}
