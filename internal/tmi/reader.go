package tmi

import "github.com/gempir/go-twitch-irc/v2"

type Reader struct {
	client         *twitch.Client
	ChanChan       chan string
	messageHandler MessageHandler
}

func NewReader() *Reader {
	r := Reader{
		client:         twitch.NewAnonymousClient(),
		ChanChan:       make(chan string),
		messageHandler: MessageHandler{},
	}

	r.client.OnPrivateMessage(r.messageHandler.handlePrivMsg)

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
