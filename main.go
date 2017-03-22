package test1

import (
	handlers "github.com/cancerballs/test1/handlers"
	utils "github.com/cancerballs/test1/utils"
	"net/http"
)

func init() {
	// Providing the buffer length as the second argument
	// makes this a buffered channel
	utils.Channel_cmd = make(chan utils.Cmd, 4)
	http.HandleFunc("/cmd", handlers.HandleRequest)
}
