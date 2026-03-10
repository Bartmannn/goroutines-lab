package server_solution

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Reviver struct {
	alivePlayers int
	currentId    int
	receiveChan  []chan Message
	requestChan  []chan Message
}

func (r *Reviver) setVerticiesChannels(
	receiveChan []chan Message,
	requestChan []chan Message,
) {
	r.receiveChan = receiveChan
	r.requestChan = requestChan
}

func (r *Reviver) nextId() string {
	value := r.currentId
	strValue := strconv.Itoa(value)
	if len(strValue) < 2 {
		return "0" + strValue
	}
	return strValue
}

func (r *Reviver) tryReviveTraveler(request chan<- Message, receive <-chan Message) bool {

	tempChannel := make(chan Message)

	request <- Message{
		requestType:  "revive",
		travelerId:   r.nextId(),
		code:         200,
		responseChan: tempChannel,
	}

	vertexResponse := <-tempChannel
	close(tempChannel)
	if vertexResponse.code >= 200 && 300 > vertexResponse.code {
		r.alivePlayers++
		r.currentId++
		return true
	}

	return false
}

func (r *Reviver) tryReviveSquatter(request chan<- Message, receive <-chan Message) bool {

	tempChannel := make(chan Message)

	request <- Message{
		requestType:  "revive",
		travelerId:   SquatterId,
		code:         200,
		responseChan: tempChannel,
	}

	vertexResponse := <-tempChannel
	close(tempChannel)
	if vertexResponse.code >= 200 && 300 > vertexResponse.code {
		return true
	}

	return false
}

func (r *Reviver) tryReviveDanger(request chan<- Message, receive <-chan Message) bool {

	tempChannel := make(chan Message)

	request <- Message{
		requestType:  "revive",
		travelerId:   DangerId,
		code:         200,
		responseChan: tempChannel,
	}

	vertexResponse := <-tempChannel
	close(tempChannel)
	if vertexResponse.code >= 200 && 300 > vertexResponse.code {
		return true
	}

	return false
}

func (r *Reviver) serviceVerticiesChannels(fromVertex <-chan Message) {
	for {
		vertexRequest := <-fromVertex
		if vertexRequest.requestType == "travelersDecrease" {
			r.alivePlayers--
		}
	}
}

func (r *Reviver) start(wg *sync.WaitGroup) {

	// traveler reviving
	wg.Add(1)
	go func() {
		for {
			time.Sleep(NewTravelerCooldown)
			randomVertex := rand.Intn(len(r.requestChan))
			if r.alivePlayers <= MaxTravelers {
				for !r.tryReviveTraveler(r.requestChan[randomVertex], r.receiveChan[randomVertex]) {
					randomVertex = rand.Intn(len(r.requestChan))
				}
			}
		}
	}()

	// squatter reviving
	wg.Add(1)
	go func() {
		for {
			time.Sleep(NewSquatterCooldown)
			randomVertex := rand.Intn(len(r.requestChan))
			for !r.tryReviveSquatter(r.requestChan[randomVertex], r.receiveChan[randomVertex]) {
				randomVertex = rand.Intn(len(r.requestChan))
			}
		}
	}()

	// danger reviving
	wg.Add(1)
	go func() {
		for {
			time.Sleep(NewDangerCooldown)
			randomVertex := rand.Intn(len(r.requestChan))
			for !r.tryReviveDanger(r.requestChan[randomVertex], r.receiveChan[randomVertex]) {
				randomVertex = rand.Intn(len(r.requestChan))
			}
		}
	}()

	for v := 0; v < len(r.receiveChan); v++ {
		wg.Add(1)
		go r.serviceVerticiesChannels(r.receiveChan[v])
	}

}
