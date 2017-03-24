package test1

import (
	"fmt"
	"testing"
)

func TestListen(t *testing.T) {
	// Test that the Listen function registers Listener("subscriber")
	// to listeners = []Listener{}
	if len(listeners) != 0 {
		t.Errorf("Expected %d, found %d", 0, len(listeners))
	}

	Listen(func(cmd Cmd) { fmt.Println(1, cmd) })

	if len(listeners) != 1 {
		t.Errorf("Expected %d, found %d", 1, len(listeners))
	}

	Listen(func(cmd Cmd) { fmt.Println(1, cmd) })
	Listen(func(cmd Cmd) { fmt.Println(1, cmd) })

	if len(listeners) != 3 {
		t.Errorf("Expected %d, found %d", 3, len(listeners))
	}
}

func TestSetupServers(t *testing.T) {
	// TODO
	// Need to figure out how this can be tested since
	// it's a big chuck of the program. Probably integration test
	// for verifying that it actually starts the http and tcp server
}
