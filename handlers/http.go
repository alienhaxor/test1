package handlers

import (
	"encoding/json"
	utils "github.com/cancerballs/test1/utils"
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

	var cmd utils.Cmd
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&cmd)
	if err != nil {
		data := map[string]string{"status": "failure"}
		jData, err := json.Marshal(data)
		if err != nil {
			// panic(err)
			return err
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jData)
		return err
	}

	defer r.Body.Close()

	// TODO: in a world where external libraries can be used one
	// could refactor this to use some sort of JSON schema
	// maybe https://github.com/xeipuuv/gojsonschema
	// ***
	// validate the struct created from the JSON
	// "" and 0 are zero values for variables declared without an explicit initial value
	// meaning a value was not supplied which makes this an invalid request.
	// Requests with Type == 0 are invalid
	if cmd.Body == "" || cmd.Type == 0 {
		// invalid JSON payload
		data := map[string]string{"status": "failure"}
		jData, err := json.Marshal(data)
		if err != nil {
			// panic(err)
			return err
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jData)
		return err
	}

	// data is valid and is sent to the channel
	utils.Channel_cmd <- cmd

	data := map[string]string{"status": "success"}
	jData, err := json.Marshal(data)
	if err != nil {
		// panic(err)
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jData)
	return
}
