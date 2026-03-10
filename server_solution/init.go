package server_solution

import (
	"fmt"
	"sync"
)

func Init() {

	var wg sync.WaitGroup

	// initiate verticies grid with basic values
	verticiesGrid := make([][]Vertex, FieldHeight)
	for y := 0; y < FieldHeight; y++ {
		verticiesGrid[y] = make([]Vertex, FieldWidth)
		for x := 0; x < FieldWidth; x++ {
			verticiesGrid[y][x] = Vertex{
				travelerId:          EmptyId,
				availableNeighbours: [4]bool{true, true, true, true},
				visiblePaths:        make([]int, 0),
				isBlocked:           false,
			}
			if x == 0 || y == 0 || x == FieldWidth-1 || y == FieldHeight-1 {
				verticiesGrid[y][x].isBlocked = true
			}
		}
	}

	// setting snapshotter channels
	snapshotReceiveChan2d := make([][]chan Message, FieldHeight)
	for y := 0; y < FieldHeight; y++ {
		snapshotReceiveChan2d[y] = make([]chan Message, FieldWidth)
		for x := 0; x < FieldWidth; x++ {
			snapshotReceiveChan2d[y][x] = make(chan Message, 1)
		}
	}

	// creating channels
	snapshotRequestsChan := make([]chan Message, 0)
	snapshotReceiveChan := make([]chan Message, 0)
	reviverReceiveChan := make([]chan Message, 0)
	reviverRequestsChan := make([]chan Message, 0)
	for y := 0; y < FieldHeight; y++ {
		for x := 0; x < FieldWidth; x++ {
			snapshotReceiveChan = append(snapshotReceiveChan, make(chan Message, 1))
			snapshotRequestsChan = append(snapshotRequestsChan, make(chan Message, 1))
			reviverReceiveChan = append(reviverReceiveChan, make(chan Message, 1))
			reviverRequestsChan = append(reviverRequestsChan, make(chan Message, 1))
		}
	}

	// adding channels to verticies
	i := 0
	for y := 0; y < FieldHeight; y++ {
		for x := 0; x < FieldWidth; x++ {
			neigboursChanRequest := make([]chan Message, 4)
			for k := 0; k < 4; k++ {
				neigboursChanRequest[k] = make(chan Message, 1)
			}
			verticiesGrid[y][x].setSnapshotChannels(snapshotRequestsChan[i], snapshotReceiveChan[i])
			verticiesGrid[y][x].setReviverChannels(reviverReceiveChan[i], reviverRequestsChan[i])
			verticiesGrid[y][x].setNeighboursRequestChannels(neigboursChanRequest)
			verticiesGrid[y][x].initNeigboursResponse()
			i++
		}
	}

	// verticies channels exchange
	for y := 0; y < len(verticiesGrid); y++ {
		for x := 0; x < len(verticiesGrid[0]); x++ {
			if y-1 >= 0 {
				verticiesGrid[y-1][x].fromNeighboursChan[2] =
					verticiesGrid[y][x].toNeighboursChan[0]
			} else {
				verticiesGrid[y][x].availableNeighbours[0] = false
				// fmt.Println(y, " ", x)
			}
			if y+1 < FieldHeight {
				verticiesGrid[y+1][x].fromNeighboursChan[0] =
					verticiesGrid[y][x].toNeighboursChan[2]
			} else {
				verticiesGrid[y][x].availableNeighbours[2] = false
				// fmt.Println(y, " ", x)
			}
			if x-1 >= 0 {
				verticiesGrid[y][x-1].fromNeighboursChan[1] =
					verticiesGrid[y][x].toNeighboursChan[3]
			} else {
				verticiesGrid[y][x].availableNeighbours[3] = false
				// fmt.Println(y, " ", x)
			}
			if x+1 < FieldWidth {
				verticiesGrid[y][x+1].fromNeighboursChan[3] =
					verticiesGrid[y][x].toNeighboursChan[1]
			} else {
				verticiesGrid[y][x].availableNeighbours[1] = false
				// fmt.Println(y, " ", x)
			}
		}
	}

	// starting verticies
	for y := 0; y < FieldHeight; y++ {
		for x := 0; x < FieldWidth; x++ {
			wg.Add(1)
			go verticiesGrid[y][x].start(&wg)
		}
	}

	// initiate snapshot server
	snapshotter := Snapshotter{}
	snapshotter.setSnapshotChannels(snapshotReceiveChan, snapshotRequestsChan)
	wg.Add(1)
	snapshotter.start(&wg)

	// initiate reviver server
	reviver := Reviver{alivePlayers: 0, currentId: 0}
	reviver.setVerticiesChannels(reviverReceiveChan, reviverRequestsChan)
	wg.Add(1)
	go reviver.start(&wg)

	fmt.Println("To stop press ENTER key")
	fmt.Scanln()

}
