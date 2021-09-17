package main_test

import (
	"testing"

	"github.com/m4tthewde/tmihooks/test"
)

func TestRun(t *testing.T) {
	t.Parallel()

	testServer := test.NewTestServer(t, test.INFINITE)

	go testServer.StartTestClient()

	for {
		select {}
	}
}
