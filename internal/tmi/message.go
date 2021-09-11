package tmi

import (
	"log"

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
		log.Println(msg.Channel, webhook)
	}
}
