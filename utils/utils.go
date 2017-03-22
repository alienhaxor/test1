package utils

type Cmd struct {
	Body string `json:"body"`
	Type int    `json:"type"`
}

// Refactor: research if feasible to save the channel to a structuct
var Channel_cmd chan Cmd
