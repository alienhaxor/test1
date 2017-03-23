package test1

import (
	"net/http"
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
