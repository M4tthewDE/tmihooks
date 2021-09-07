package webhook

type Webhook struct {
	Channels []string `json:"channels"`
	URI      string   `json:"uri"`
	Nonce    string   `json:"nonce"`
	Status   string
}
