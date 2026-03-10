package server_solution

import "time"

type Message struct {
	requestType  string
	travelerId   string
	direction    int
	code         int
	dirs         []int
	responseChan chan Message
	restTime     time.Duration
}
