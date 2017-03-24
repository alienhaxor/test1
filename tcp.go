package test1

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

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
			conn.Write([]byte("Json error."))
			return
		}

		_, err = validateData(cmd)
		if err != nil {
			conn.Write([]byte(`{"status": "failure"}`))
			return
		}

		channelCmd <- cmd

		data := []byte(`{"status": "success"}`)
		conn.Write([]byte(data))
	}
}
