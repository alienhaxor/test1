package test1

import (
	"encoding/json"
	"net/http"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		CmdHandler(w, r)
	}
}

func CmdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cmd Cmd
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(&cmd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "failure"}`))
		return
	}

	_, err = validateData(cmd)
	if err != nil {
		// invalid JSON payload
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"status": "failure"}`))
		return
	}

	// data is valid and is sent to the channel
	channelCmd <- cmd

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
	return
}
