package test1

import (
	"encoding/json"
	"net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "POST":
		err = CmdHandler(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CmdHandler(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "application/json")

	var cmd Cmd
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err = decoder.Decode(&cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "failure"}`))
		return err
	}

	// TODO: in a world where external libraries can be used one
	// could refactor this to use some sort of JSON schema
	// maybe https://github.com/xeipuuv/gojsonschema
	// TODO: refactor validate and send to channel to 1 function that is shared
	// by http CmdHandler and tcp handleConnection
	// ***
	// validate the struct created from the JSON
	// "" and 0 are zero values for variables declared without an explicit initial value
	// meaning a value was not supplied which makes this an invalid request.
	// Requests with Type == 0 are invalid
	if cmd.Body == "" || cmd.Type == 0 {
		// invalid JSON payload
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "failure"}`))
		return err
	}

	// data is valid and is sent to the channel
	channelCmd <- cmd

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
	return
}
