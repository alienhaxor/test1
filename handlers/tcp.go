package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	utils "github.com/cancerballs/test1/utils"
	"net"
	"time"
)

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
		var cmd utils.Cmd
		err2 := json.Unmarshal(bytes, &cmd)
		if err2 != nil {
			fmt.Println(err)
			conn.Write([]byte("Json error."))
			return
		}

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
			data := []byte(`{"status": "failure"}`)
			conn.Write([]byte(data))
			return
		}

		utils.Channel_cmd <- cmd
		data := []byte(`{"status": "success"}`)
		conn.Write([]byte(data))
	}
}

func TcpServer() {
	listener, err := net.Listen("tcp", ":8888")
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
