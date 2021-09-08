package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/m4tthewde/tmihooks/internal/config"
	"github.com/m4tthewde/tmihooks/internal/db"
	"github.com/m4tthewde/tmihooks/internal/structs"
	"github.com/m4tthewde/tmihooks/internal/util"
)

type RouteHandler struct {
	dbHandler *db.DatabaseHandler
}

func NewRouteHandler(config *config.Config) *RouteHandler {
	dbHandler := &db.DatabaseHandler{
		Config: config,
	}

	return &RouteHandler{
		dbHandler: dbHandler,
	}
}

// register new webhook.
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

// get all webhooks.
func (rh *RouteHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

// get all webhooks.
func (rh *RouteHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

	req, err := http.NewRequestWithContext(ctx, "POST", webhook.URI, bytes.NewBuffer(confirmationJSON))
	if err != nil {
		log.Println(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK || string(body) != confirmation.Challenge {
		log.Println("webhook wasn't confirmed properly!")
		rh.dbHandler.Delete(webhook)
	} else {
		n := rh.dbHandler.SetConfirmed(confirmation.ID)
		if n != 1 {
			panic("not exaclty one webhook was confirmed")
		}
	}
}
