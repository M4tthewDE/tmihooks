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
	"github.com/m4tthewde/tmihooks/internal/structs"
)

type MessageHandler struct {
	dbHandler   *db.DatabaseHandler
	webhooks    []structs.Webhook
	WebhookChan chan structs.Webhook
}

func NewMessageHandler(config *config.Config) *MessageHandler {
	mh := MessageHandler{
		dbHandler: &db.DatabaseHandler{
			Config: config,
		},
		webhooks:    make([]structs.Webhook, 0, 10),
		WebhookChan: make(chan structs.Webhook),
	}

	return &mh
}

func (mh *MessageHandler) handlePrivMsg(msg twitch.PrivateMessage) {
	// check in slice for webhooks with channel
	// if that is too slow too, a map could be used for better performance
	webhooks := mh.GetWebhooksWithChannel(msg.Channel)

	for _, webhook := range webhooks {
		mh.sendToEndpoint(msg, webhook.URI)
	}
}

func (mh *MessageHandler) WebhookListener() {
	for webhook := range mh.WebhookChan {
		mh.webhooks = append(mh.webhooks, webhook)

		for _, wh := range mh.webhooks {
			log.Println(wh.Channels)
		}
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

func (mh *MessageHandler) GetWebhooksWithChannel(channel string) []structs.Webhook {
	result := make([]structs.Webhook, 0)

	for _, webhook := range mh.webhooks {
		for _, c := range webhook.Channels {
			if c == channel {
				result = append(result, webhook)
			}
		}
	}
	return result
}
