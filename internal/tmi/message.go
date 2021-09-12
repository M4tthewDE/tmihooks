package tmi

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/db"
)

type MessageHandler struct {
	dbHandler *db.DatabaseHandler
}

func NewMessageHandler(config *config.Config) *MessageHandler {
	mh := MessageHandler{
		dbHandler: &db.DatabaseHandler{
			Config: config,
		},
	}

	return &mh
}

func (mh *MessageHandler) handlePrivMsg(msg twitch.PrivateMessage) {
	// check in slice for webhooks with channel
	// if that is too slow too, a map could be used for better performance
	uris, err := db.GetURIs(msg.Channel)
	if err != nil {
		panic(err)
	}

	for _, uri := range uris {
		mh.sendToEndpoint(msg, uri)
	}
}

func (mh *MessageHandler) sendToEndpoint(msg twitch.PrivateMessage, uri string) {
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer(msgJSON))
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp == nil {
		// TODO do something in this case?
		log.Println("no response")
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO their problem I guess
		log.Println(resp.StatusCode)
	}
}
