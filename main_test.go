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

func TestValidateData(t *testing.T) {
	validData := Cmd{"King of Rock", 666}
	invalidData := Cmd{"", 666}
	invalidData2 := Cmd{"", 0}
	invalidData4 := Cmd{}

	_, err := validateData(validData)
	if err != nil {
		t.Errorf("This data should be valid")
	}

	_, err = validateData(invalidData)
	if err == nil {
		t.Errorf("This data is invalid")
	}

	_, err = validateData(invalidData2)
	if err == nil {
		t.Errorf("This data is invalid")
	}

	_, err = validateData(invalidData4)
	if err == nil {
		t.Errorf("This data is invalid")
	}
}
