package tmi

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/m4tthewde/tmihooks/internal/db"
)

type MessageHandler struct {
	dbHandler *db.DatabaseHandler
}

func (mh *MessageHandler) handlePrivMsg(msg twitch.PrivateMessage) {
	// this could have terrible performance, but it works for now
	// if performance becomes critical, webhooks can be kept in memory
	webhooks, err := mh.dbHandler.GetWebhooksByChannel(msg.Channel)
	if err != nil {
		log.Panic(err)
	}
	for _, webhook := range webhooks {
		// send msg to endpoint
		mh.sendToEndpoint(msg, webhook.URI)
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
		log.Println("message wasn't received properly")
		// TODO do something in this case?
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("message wasn't received properly!")
		// TODO their problem I guess
	}
}
