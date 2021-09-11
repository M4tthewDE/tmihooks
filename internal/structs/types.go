package structs

type Webhook struct {
	Channels    []string `json:"channels"`
	URI         string   `json:"uri"`
	RegisterURI string   `json:"register_uri"`
	Nonce       string   `json:"nonce"`
	Status      string
}

type Confirmation struct {
	Nonce     string
	ID        string
	Challenge string
}
