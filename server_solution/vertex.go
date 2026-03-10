package server_solution

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Vertex struct {
	travelerId              string
	lastDirection           int
	availableNeighbours     [4]bool
	visiblePaths            []int
	isBlocked               bool
	liveTime                SecondsTimer
	sendToSnapshotChan      chan Message
	receiveFromSnapshotChan chan Message
	sendToReviverChan       chan Message
	gainFromReviverChan     chan Message
	fromNeighboursChan      []chan Message
	toNeighboursChan        []chan Message
}

func directionName(direction int) string {
	switch direction {
	case DirectionUp:
		return "up"
	case DirectionRight:
		return "right"
	case DirectionDown:
		return "down"
	case DirectionLeft:
		return "left"
	default:
		return "unknown"
	}
}

func logEvent(format string, args ...any) {
	if AreCommentLabels {
		fmt.Printf(format+"\n", args...)
	}
}

func (v *Vertex) NewSecondsTimer(t time.Duration) *SecondsTimer {
	return &SecondsTimer{time.NewTimer(t), time.Now().Add(t)}
}

func (v *Vertex) setSnapshotChannels(receiveFromSnapshotChan chan Message, sendToSnapshotChan chan Message) {
	v.receiveFromSnapshotChan = receiveFromSnapshotChan
	v.sendToSnapshotChan = sendToSnapshotChan
}

func (v *Vertex) setReviverChannels(sendToReviverChan chan Message, gainFromReviverChan chan Message) {
	v.sendToReviverChan = sendToReviverChan
	v.gainFromReviverChan = gainFromReviverChan
}

func (v *Vertex) setNeighboursRequestChannels(toNeighboursChan []chan Message) {
	v.toNeighboursChan = toNeighboursChan
}

func (v *Vertex) initNeigboursResponse() {
	v.fromNeighboursChan = make([]chan Message, 4)
}

func (v *Vertex) revive(id string) {
	v.travelerId = id
	v.lastDirection = DirectionNone
	if len(id) == 2 {
		v.notifyBlock()
	}
	v.updateSnapshot("reviveUpdate")
}

func (v *Vertex) notifyBlock() {
	for i := 0; i < len(v.toNeighboursChan); i++ {
		if v.toNeighboursChan[i] != nil {
			v.toNeighboursChan[i] <- Message{requestType: "block", code: 200}
		}
	}
}

func (v *Vertex) notifyUnblock() {
	for i := 0; i < len(v.toNeighboursChan); i++ {
		if v.toNeighboursChan[i] != nil {
			v.toNeighboursChan[i] <- Message{requestType: "unblock", code: 200}
		}
	}
}

func (v *Vertex) updateSnapshot(requestType string) {
	v.sendToSnapshotChan <- Message{
		requestType: requestType,
		travelerId:  v.travelerId,
		direction:   v.lastDirection,
		code:        200,
		dirs:        v.visiblePaths,
	}
	v.visiblePaths = make([]int, 0)
}

func (v *Vertex) serviceTimer() {
	<-v.liveTime.timer.C
	switch v.travelerId {
	case SquatterId:
		logEvent("Squatter disappeared.")
	case DangerId:
		logEvent("Danger disappeared.")
	default:
		logEvent("Traveler %s disappeared.", v.travelerId)
	}
	v.travelerId = EmptyId
	v.lastDirection = DirectionNone
	v.updateSnapshot("liveTimeOut")
}

func (v *Vertex) serviceReviver(toReviver chan<- Message, fromReviver <-chan Message) {
	for {
		select {
		case reviverRequest := <-fromReviver:
			if v.travelerId != "x" || v.isBlocked {
				// toReviver <- Message{code: 500}
				reviverRequest.responseChan <- Message{code: 500}
			} else {
				v.revive(reviverRequest.travelerId)
				if v.travelerId == SquatterId {
					logEvent("A squatter appeared.")
					v.liveTime = *v.NewSecondsTimer(SquatterLiveTime)
					go v.serviceTimer()
				} else if v.travelerId == DangerId {
					logEvent("A danger appeared.")
					v.liveTime = *v.NewSecondsTimer(DangerLiveTime)
					go v.serviceTimer()
				} else {
					logEvent("Traveler %s entered the board.", v.travelerId)
				}
				// toReviver <- Message{code: 200}
				reviverRequest.responseChan <- Message{code: 200}
			}
		}
	}
}

func (v *Vertex) tryMove() {
	movements := make([]int, 0)
	for w := 0; w < len(v.availableNeighbours); w++ {
		if v.availableNeighbours[w] {
			movements = append(movements, w)
		}
	}

	if len(movements) > 0 {
		tempRespone := make(chan Message)
		moveDir := movements[rand.Intn(len(movements))]
		message := Message{
			requestType:  "movement",
			travelerId:   v.travelerId,
			direction:    moveDir,
			code:         200,
			responseChan: tempRespone,
		}

		if v.travelerId == SquatterId {
			v.liveTime.Stop()
			message.restTime = v.liveTime.TimeRemaining()
		}

		v.toNeighboursChan[moveDir] <- message

		select {
		case neigbourResponse := <-tempRespone:
			if neigbourResponse.code >= 200 && 300 > neigbourResponse.code {
				logEvent("Traveler %s moved %s.", v.travelerId, directionName(moveDir))
				v.travelerId = EmptyId
				v.lastDirection = DirectionNone
				v.notifyUnblock()
				v.visiblePaths = append(v.visiblePaths, moveDir)
				v.updateSnapshot("moveOutUpdate")
			} else if neigbourResponse.code >= 500 {
				if neigbourResponse.requestType == "kill" {
					v.travelerId = EmptyId
					v.lastDirection = DirectionNone
					v.notifyUnblock()
				} else {
					v.tryMove()
				}
			}
		}

	}
}

func (v *Vertex) serviceNeigbour(dir int, toNeighbour chan<- Message, fromNeigbour <-chan Message) {
	for {
		select {
		case neigbourRequest := <-fromNeigbour:
			switch neigbourRequest.requestType {
			case "block":
				v.availableNeighbours[dir] = false
				toNeighbour <- Message{code: 200}
			case "unblock":
				v.availableNeighbours[dir] = true
				toNeighbour <- Message{code: 200}
			case "movement":
				if v.travelerId == EmptyId && !v.isBlocked {
					v.travelerId = neigbourRequest.travelerId
					v.lastDirection = neigbourRequest.direction
					if v.travelerId == SquatterId {
						v.liveTime = *v.NewSecondsTimer(neigbourRequest.restTime)
						go v.serviceTimer()
					} else {
						v.notifyBlock()
					}
					v.updateSnapshot("moveToUpdate")
					neigbourRequest.responseChan <- Message{code: 200}
				} else if v.travelerId == SquatterId {
					v.tryMove()
					v.travelerId = neigbourRequest.travelerId
					v.lastDirection = neigbourRequest.direction
					logEvent("Traveler %s displaced a squatter.", v.travelerId)
					v.notifyBlock()
					v.updateSnapshot("moveToUpdate")
					neigbourRequest.responseChan <- Message{code: 200}
				} else if v.travelerId == DangerId {
					logEvent("Traveler %s stepped on danger and disappeared.", neigbourRequest.travelerId)
					v.travelerId = EmptyId
					v.lastDirection = DirectionNone
					v.liveTime.Stop()
					v.updateSnapshot("moveToUpdate")
					neigbourRequest.responseChan <- Message{code: 200, requestType: "kill"}
					v.sendToReviverChan <- Message{requestType: "travelersDecrease"}
				} else {
					neigbourRequest.responseChan <- Message{code: 500}
				}
			}
		}
	}
}

func (v *Vertex) start(wg *sync.WaitGroup) {

	wg.Add(1)
	go func() {
		for {
			cooldown := rand.Intn(MovementCooldownSup-MovementCooldownInf) + MovementCooldownInf
			time.Sleep(time.Millisecond * time.Duration(cooldown))
			if v.travelerId == SquatterId {
				// do nothing for a while
			} else if v.travelerId == DangerId {
				// do nothing for a while
			} else if v.travelerId == EmptyId {
				// do nothing for a while
			} else {
				v.tryMove()
			}
		}
	}()

	wg.Add(1)
	go v.serviceReviver(v.sendToReviverChan, v.gainFromReviverChan)

	for n := 0; n < len(v.availableNeighbours); n++ {
		wg.Add(1)
		go v.serviceNeigbour(n, v.toNeighboursChan[n], v.fromNeighboursChan[n])
	}
}
