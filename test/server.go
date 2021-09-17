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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	t             *testing.T
	config        *config.Config
	server        *http.Server
	router        *chi.Mux
	webhook       *structs.Webhook
	StopChan      chan int
	interruptChan chan os.Signal
	testType      Type
}

func NewTestServer(t *testing.T, testType Type) *Server {
	config := config.GetConfig("../test_config.yml")

	ts := Server{
		t:             t,
		config:        config,
		webhook:       nil,
		server:        &http.Server{Addr: ":7070"},
		router:        chi.NewRouter(),
		StopChan:      make(chan int, 1),
		interruptChan: make(chan os.Signal, 1),
		testType:      testType,
	}
	signal.Notify(ts.interruptChan, os.Interrupt)

	return &ts
}

func (ts *Server) StartTestClient() {
	ts.router.Use(middleware.Logger)
	ts.router.Post("/register", ts.Register())
	ts.router.Post("/chat", ts.chat())

	ts.server.Handler = ts.router

	go func() {
		err := ts.server.ListenAndServe()
		if err != nil {
			log.Println("server closed")
		}
	}()

	<-ts.interruptChan
	log.Println("stopping main server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := ts.server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("done")
	defer cancel()
	if err := ts.server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("done")
}

func (ts *Server) Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var confirmation structs.Confirmation

		err := json.NewDecoder(r.Body).Decode(&confirmation)
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
}

func (ts *Server) chat() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg twitch.PrivateMessage

		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			panic(err)
		}

		assert.Equal(ts.t, "tmiloadtesting2", msg.Channel)

		if ts.testType == REGISTER {
			log.Println("received message, attempting graceful shutdown")

			// shutdown main server
			ts.ShutdownMainServer()

			// shutdown test server
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				panic(err)
			}

			err = p.Signal(os.Interrupt)
			if err != nil {
				panic(err)
			}
			ts.StopChan <- 0
		}
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
