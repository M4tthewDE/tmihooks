package api

import (
	"encoding/json"
	"net/http"

	"github.com/m4tthewde/tmihooks/internal/db"
	"github.com/m4tthewde/tmihooks/internal/webhook"
)

type RouteHandler struct {
}

// register new webhook
func (rh *RouteHandler) Register() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var webhook webhook.Webhook
		json.NewDecoder(r.Body).Decode(&webhook)
		webhook.Status = "unconfirmed"

		db.Insert(webhook)
	}
}

// get all webhooks
func (rh *RouteHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

// get all webhooks
func (rh *RouteHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
