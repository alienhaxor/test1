package test1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCmdHandler(t *testing.T) {
	// TODO:
	// The channel is tightly coupled to the httpHandler hence
	// it has to be tested here, otherwise the test hangs indefinitely
	// ***
	// need to create the channel CmdHandler uses
	channelCmd = make(chan Cmd, 4)

	mux := http.NewServeMux()
	mux.HandleFunc("/cmd", HandleRequest)

	writer := httptest.NewRecorder()
	json := strings.NewReader(`{"body": "King of Rock", "type": 666}`)
	request, _ := http.NewRequest("POST", "/cmd", json)
	mux.ServeHTTP(writer, request)

	// test channel(internal purposes)
	expected := Cmd{"King of Rock", 666}
	found := <-channelCmd
	if found != expected {
		t.Errorf("Expected %s, found %s", expected, found)
	}

	// Check the http response code
	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	// Check the response body
	expected_result := `{"status": "success"}`
	if writer.Body.String() != expected_result {
		t.Errorf("handler returned unexpected body: got %v want %v",
			writer.Body.String(), expected_result)
	}

	writer2 := httptest.NewRecorder()
	// invalid json, should fail
	json2 := strings.NewReader(`{"body2": "King of Rock", "type": 666}`)
	request2, _ := http.NewRequest("POST", "/cmd", json2)
	mux.ServeHTTP(writer2, request2)

	// expecting bad request since invalid json
	if writer2.Code != 400 {
		t.Errorf("Response code is %v", writer2.Code)
	}

	// expecting failure
	expected_result2 := `{"status": "failure"}`
	if writer2.Body.String() != expected_result2 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			writer2.Body.String(), expected_result2)
	}
}
