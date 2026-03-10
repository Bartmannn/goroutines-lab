package server_solution

import "time"

type Message struct {
	requestType  string
	travelerId   string
	code         int
	dirs         []int
	responseChan chan Message
	restTime     time.Duration
}
