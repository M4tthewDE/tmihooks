package main_test

import (
	"testing"

	"github.com/m4tthewde/tmihooks/test"
)

func TestRunApplication(t *testing.T) {
	t.Parallel()

	testServer := test.NewTestServer(t)

	go testServer.StartTestClient()

	for {
		select {}
	}
}
