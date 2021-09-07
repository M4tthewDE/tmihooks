package main

import "github.com/m4tthewde/tmihooks/internal/api"

func main() {
	server := api.NewServer("1500")
	server.Run()
}
