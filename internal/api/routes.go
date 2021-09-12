package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/db"
	"github.com/m4tthewde/tmihooks/internal/structs"
	"github.com/m4tthewde/tmihooks/internal/tmi"
	"github.com/m4tthewde/tmihooks/internal/util"
)

type RouteHandler struct {
	dbHandler *db.DatabaseHandler
	reader    *tmi.Reader
}

func NewRouteHandler(config *config.Config, reader *tmi.Reader) *RouteHandler {
	dbHandler := &db.DatabaseHandler{
		Config: config,
	}

	webhooks, err := dbHandler.GetAllWebhooks()
	if err != nil {
		panic(err)
	}

	for _, webhook := range webhooks {
		db.AddWebhook(webhook)
		for _, channel := range webhook.Channels {
			reader.ChanChan <- channel
		}
	}

	return &RouteHandler{
		dbHandler: dbHandler,
		reader:    reader,
	}
}

// register new webhook.
// /register
func (rh *RouteHandler) Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var webhook structs.Webhook

		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			panic(err)
		}

		webhook.Status = "unconfirmed"

		id := rh.dbHandler.Insert(&webhook)
		confirmation := structs.Confirmation{
			Nonce:     webhook.Nonce,
			ID:        id.Hex(),
			Challenge: util.RandomString(32),
		}

		rh.ConfirmWebhook(&confirmation, &webhook)
	}
}

// get webhook by id.
// /get?id=asdf1234
func (rh *RouteHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.URL.Query()["id"]
		if !ok {
			http.Error(w, "Bad id.", http.StatusBadRequest)
		}

		webhook, err := rh.dbHandler.GetWebhook(id[0])
		if err != nil {
			http.Error(w, "Webhook not found.", http.StatusBadRequest)
		}

		jsonWebhook, err := json.Marshal(webhook)
		if err != nil {
			panic(err)
		}

		_, err = w.Write(jsonWebhook)
		if err != nil {
			panic(err)
		}
	}
}

// delete webhook by id.
// /delete
func (rh *RouteHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.URL.Query()["id"]
		if !ok {
			http.Error(w, "Bad id.", http.StatusBadRequest)
		}

		n, err := rh.dbHandler.DeleteByID(id[0])
		if err != nil {
			http.Error(w, "Webhook not found.", http.StatusBadRequest)
			return
		}
		if n == 0 {
			http.Error(w, "Webhook not found.", http.StatusBadRequest)
		}

		webhook, err := rh.dbHandler.GetWebhook(id[0])
		if err != nil {
			panic(err)
		}
		db.DeleteWebhook(webhook)

		_, err = w.Write([]byte(id[0]))
		if err != nil {
			panic(err)
		}
	}
}

func (rh *RouteHandler) ConfirmWebhook(confirmation *structs.Confirmation, webhook *structs.Webhook) {
	confirmationJSON, err := json.Marshal(confirmation)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", webhook.RegisterURI, bytes.NewBuffer(confirmationJSON))
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	if resp == nil {
		log.Println("webhook wasn't confirmed properly!")
		rh.dbHandler.Delete(webhook)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusOK || string(body) != confirmation.Challenge {
		log.Println("webhook wasn't confirmed properly!")
		rh.dbHandler.Delete(webhook)
	} else {
		n := rh.dbHandler.SetConfirmed(confirmation.ID)
		if n != 1 {
			panic("not exaclty one webhook was confirmed")
		} else {
			// webhook was confirmed successfully.
			db.AddWebhook(webhook)
			for _, channel := range webhook.Channels {
				rh.reader.ChanChan <- channel
			}
		}
	}
}

func (rh *RouteHandler) Shutdown() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := os.FindProcess(os.Getpid())
		if err != nil {
			panic(err)
		}

		err = p.Signal(os.Interrupt)
		if err != nil {
			panic(err)
		}
	}
}
