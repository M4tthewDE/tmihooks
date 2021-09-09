package tmi

import (
	"log"

	"github.com/gempir/go-twitch-irc/v2"
)

type MessageHandler struct {
}

func (mh *MessageHandler) handlePrivMsg(msg twitch.PrivateMessage) {
	log.Println(msg.Channel, msg.Message)
}
