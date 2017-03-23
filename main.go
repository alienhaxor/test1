package test1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Cmd struct {
	Body string `json:"body"`
	Type int    `json:"type"`
}

func init() {
	// Providing the buffer length as the second argument
	// makes this a buffered channel
	channelCmd = make(chan Cmd, 4)
	exit = make(chan bool)
}

// Refactor: research if feasible to save the channel to a structuct
var channelCmd chan Cmd
var exit chan bool
var listeners = []Listener{}

type Listener func(cmd Cmd)

func Listen(listener Listener) {
	listeners = append(listeners, listener)
}

func SetupServers(httpPort string, tcpPort string) {
	http.HandleFunc("/cmd", HandleRequest)
	go TcpServer(tcpPort)

	go func() {
		for {
			select {
			case in := <-channelCmd:
				for _, l := range listeners {
					go l(in)
				}
			case <-exit:
				return
			}
		}
	}()

	http.ListenAndServe(httpPort, nil)
}

func TcpServer(tcpPort string) {
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Close connection when this function ends
	defer func() {
		conn.Close()
	}()

	timeoutDuration := 5 * time.Second
	bufReader := bufio.NewReader(conn)

	for {
		// Set a deadline for reading. Read operation will fail if no data
		// is received after deadline.
		conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		// Read tokens delimited by newline
		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		var cmd Cmd
		err2 := json.Unmarshal(bytes, &cmd)
		if err2 != nil {
			fmt.Println(err)
			conn.Write([]byte("Json error."))
			return
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
			conn.Write([]byte(`{"status": "failure"}`))
			return
		}

		channelCmd <- cmd
		data := []byte(`{"status": "success"}`)
		conn.Write([]byte(data))
	}
}

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
