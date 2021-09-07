package webhook

import "github.com/google/uuid"

type Webhook struct {
	UUID     uuid.UUID
	Channels []string `json:"channels"`
	URI      string   `json:"uri"`
	Nonce    string   `json:"nonce"`
	Status   string
}
