package server_solution

import (
	"fmt"
	"sync"
	"time"
)

type Field struct {
	id        string
	isVisible bool
}

type Snapshotter struct {
	board             [][]Field
	fromVerticiesChan []chan Message
	toVerticiesChan   []chan Message
}

func (s *Snapshotter) draw() {
	currWidth := FieldWidth*2 - FieldWidth/2 - 1
	for i := 0; i < currWidth; i++ {
		fmt.Print("_-")
	}
	fmt.Println()
	for y := 0; y < len(s.board); y++ {
		for x := 0; x < len(s.board[0]); x++ {
			if s.board[y][x].isVisible {
				fmt.Print(s.board[y][x].id)
				if s.board[y][x].id == HorizontalPath || s.board[y][x].id == VerticalPath {
					s.board[y][x].isVisible = false
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	for i := 0; i < currWidth; i++ {
		fmt.Print("_-")
	}
	fmt.Println()
}

func (s *Snapshotter) setSnapshotChannels(fromVerticiesChan []chan Message, toVerticiesChan []chan Message) {
	s.fromVerticiesChan = fromVerticiesChan
	s.toVerticiesChan = toVerticiesChan
}

func (s *Snapshotter) buildBoard(wg *sync.WaitGroup) {
	s.board = make([][]Field, FieldHeight*2-1)
	i := 0
	for y := 0; y < FieldHeight*2-1; y++ {
		s.board[y] = make([]Field, FieldWidth*2-1)
		for x := 0; x < FieldWidth*2-1; x++ {
			s.board[y][x] = Field{id: "  ", isVisible: true}
			if x%2 == 0 && y%2 == 0 {
				s.board[y][x].isVisible = true
				wg.Add(1)
				go s.listenVertex(x, y, s.toVerticiesChan[i], s.fromVerticiesChan[i])
				i++
			} else if x%2 == 1 && y%2 == 0 {
				s.board[y][x].id = HorizontalPath
			} else if x%2 == 0 && y%2 == 1 {
				s.board[y][x].id = VerticalPath
			}
		}
	}
	// fmt.Println(len(s.board), "x", len(s.board[0]))
}

func (s *Snapshotter) listenVertex(x, y int, toVertex chan<- Message, fromVertex <-chan Message) {
	for {
		select {
		case vertexResponse := <-fromVertex:
			// fmt.Println(x, "x", y, " | ", vertexResponse.requestType, " ", vertexResponse.travelerId, " ", vertexResponse.code)
			if vertexResponse.code >= 200 && 300 > vertexResponse.code {
				// if x == 0 || y == 0 || x == len(s.board[0])-1 || y == len(s.board)-1 {
				// 	fmt.Println("Check: ", vertexResponse.requestType)
				// }
				switch vertexResponse.travelerId {
				case EmptyId:
					s.board[y][x].id = "  "
				case SquatterId:
					s.board[y][x].id = "**"
				case DangerId:
					s.board[y][x].id = "##"
				default:
					s.board[y][x].id = vertexResponse.travelerId
				}
				for _, dir := range vertexResponse.dirs {
					switch dir {
					case 0:
						s.board[y-1][x].isVisible = true
					case 1:
						s.board[y][x+1].isVisible = true
					case 2:
						s.board[y+1][x].isVisible = true
					case 3:
						s.board[y][x-1].isVisible = true
					}
				}
			}
		}
	}

}

func (s *Snapshotter) updaterBoard(wg *sync.WaitGroup) {
	var x, y int = 0, 0
	for v := 0; v < len(s.toVerticiesChan); v++ {
		x = (v % FieldWidth) * 2
		if v%FieldWidth == 0 && v > 0 {
			y += 2
		}

		wg.Add(1)
		// fmt.Println(x, " ", y, " ", v, " ")
		go s.listenVertex(x, y, s.toVerticiesChan[v], s.fromVerticiesChan[v])
	}
}

func (s *Snapshotter) start(wg *sync.WaitGroup) {
	s.buildBoard(wg)
	// s.updaterBoard(wg)

	wg.Add(1)
	go func() {
		for {
			time.Sleep(BoardRefreshRate)
			s.draw()
		}
	}()
}
