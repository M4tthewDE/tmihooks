package tmi

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/m4tthewde/tmihooks/internal/config"
)

type Reader struct {
	client         *twitch.Client
	ChanChan       chan string
	MessageHandler MessageHandler
}

func NewReader(config *config.Config) *Reader {
	r := Reader{
		client:         twitch.NewAnonymousClient(),
		ChanChan:       make(chan string),
		MessageHandler: *NewMessageHandler(config),
	}

	r.client.OnPrivateMessage(r.MessageHandler.handlePrivMsg)

	go r.MessageHandler.WebhookListener()

	return &r
}

func (r *Reader) Read() {
	go r.ChannelListener()

	err := r.client.Connect()
	if err != nil {
		panic(err)
	}
}

func (r *Reader) ChannelListener() {
	for channel := range r.ChanChan {
		r.client.Join(channel)
	}
}
