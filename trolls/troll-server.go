package trolls

import (
        "log"
        "net/http"
        "strconv"
        "math"

        "code.google.com/p/go.net/websocket"
)

const NEW_CONNECTION_ENDPOINT = "/connect";



// troll server
type Server struct {
        trolls    		map[int]*Troll

        gridItemsMap	map[int]*GridItem
        updateMap		map[int]*GridItem
        grid			[][]int // (x,y) position mapped to key of GridItem in gridItemsMap

        addCh     		chan *Troll
        delCh     		chan *Troll
        messageCh 		chan *IncomingMessage
        doneCh    		chan bool
        errCh     		chan error
}

// Create new troll server.
func NewServer() *Server {
	trolls := make(map[int]*Troll)
	gridItemsMap := make(map[int]*GridItem)
	updateMap := make(map[int]*GridItem)

	grid := NewGrid()

	addCh := make(chan *Troll)
	delCh := make(chan *Troll)
	messageCh := make(chan *IncomingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	s :=  &Server{
		trolls,
		gridItemsMap,
		updateMap,
		grid,
		addCh,
		delCh,
		messageCh,
		doneCh,
		errCh,
	}
	// add the food button to the grid
	foodButton := NewFoodButton()
	s.gridItemsMap[FOODBUTTON_ID] = foodButton

	return s
}


// called by troll client when client recieves message -- then recieveMessage called
func (s *Server) RecieveMessage(msg *IncomingMessage) {
	s.messageCh <- msg
}
func (s *Server) AddTrollConnection(t *Troll) {
        s.addCh <- t
}
func (s *Server) Del(t *Troll) {
        s.delCh <- t
}
func (s *Server) Done() {
        s.doneCh <- true
}
func (s *Server) Err(err error) {
        s.errCh <- err
}
func (s *Server) sendAll(msg *OutgoingMessage) {
	for _, t := range s.trolls {
		t.Write(msg)
	}
}

func (s *Server) sendUpdateMessage() {
	var msg *OutgoingMessage
	msg = OutgoingUpdateMessage(0, s.updateMap)
	s.sendAll(msg)

	// clear out updateMap
	s.updateMap = make(map[int]*GridItem)
}
func (s *Server) sendItemsMessage(trollID int) {
	var msg *OutgoingMessage
	msg = OutgoingItemsMessage(trollID, s.gridItemsMap)
	s.trolls[trollID].Write(msg)
}

func (s *Server) recievePingMessage(trollID int) {
	var msg *OutgoingMessage
	msg = OutgoingPingMessage(trollID)
	s.trolls[trollID].Write(msg)
}
func (s *Server) recieveItemsMessage(trollID int) {
	s.sendItemsMessage(trollID)
}

func (s *Server) recieveMessageMessage(trollID int, data map[string]string) {
	log.Println("TODO: recieveMessageMessage")
}
func (s *Server) recieveMoveMessage(trollID int, data map[string]string) {
	// extract the troll client's move from the message data
	moveX, _ := strconv.Atoi(data["x"])
	moveY, _ := strconv.Atoi(data["y"])
	// retrieve troll client's current position
	currentX := s.gridItemsMap[trollID].Coordinates["x"]
	currentY := s.gridItemsMap[trollID].Coordinates["y"]
	// calculate requested new position coordinates
	requestedX := (currentX + moveX)
	requestedY := (currentY + moveY)

	// collision detection with grid boundaries
	if (requestedX < 0 || requestedX >= GRID_WIDTH || requestedY < 0 || requestedY >= GRID_HEIGHT) {
		s.sendItemsMessage(trollID)
		return
	}
	// collision detection with other trolls
	if (s.grid[requestedX][requestedY] != 0) { 
		s.sendItemsMessage(trollID)
		return
	} else {
		// move that troll
		s.grid[currentX][currentY] = 0
		s.grid[requestedX][requestedY] = trollID

		s.gridItemsMap[trollID].Coordinates["x"] = requestedX
		s.gridItemsMap[trollID].Coordinates["y"] = requestedY

		s.updateMap[trollID] = s.gridItemsMap[trollID]
		s.sendUpdateMessage()
	}	
}


// when troll client recieves message, sends the IncomingMessage to server to be handled
func (s *Server) recieveMessage(msg *IncomingMessage) {
	log.Println("incoming msg: ", msg)

	switch msg.Type {
	case "ping":
		s.recievePingMessage(msg.LocalTroll)
	case "items":
		s.recieveItemsMessage(msg.LocalTroll)
	case "message":
		s.recieveMessageMessage(msg.LocalTroll, msg.Data)
	case "move":
		s.recieveMoveMessage(msg.LocalTroll, msg.Data)
	default:
		log.Println("Unknown message type recieved: ", msg.Type)
	}
}
func (s *Server) addTrollConnection(t *Troll) {
	log.Println("addTrollConnection *****")

	gi := t.NewGridItem()
	// find a cell for the new troll
	x := gi.Coordinates["x"]
	for (s.grid[x][0] != 0) {
		x = int(math.Mod(float64(x + 1), 9))
	}
	s.grid[x][0] = t.id
	gi.Coordinates["x"] = x


	s.updateMap[t.id] = gi
	s.sendUpdateMessage()

	s.trolls[t.id] = t
	s.gridItemsMap[t.id] = gi
	log.Println("Added new troll - Now", len(s.trolls), "trolls connected.")
	s.sendItemsMessage(t.id)
}
func (s *Server) deleteTrollConnection(t *Troll) {
	gi := s.gridItemsMap[t.id]

	// set troll to be deleted in updateMap
	gi.Name = "DELETE"
	s.updateMap[t.id] = gi

	// delete troll
	RemoveGridItem(s.grid, gi)
	delete(s.trolls, t.id)
	delete(s.gridItemsMap, t.id)

	// send update message
	s.sendUpdateMessage()
}

// Listen and serve - serves client connection and broadcast request.
func (s *Server) Listen() {
	log.Println("Troll server listening...")

	// websocket handler
	onConnect := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		troll := NewTroll(ws, s)
		s.AddTrollConnection(troll)
		troll.Listen()
	}
	http.Handle(NEW_CONNECTION_ENDPOINT, websocket.Handler(onConnect))

	for {
		select {

			// Add new a client
			case t := <-s.addCh:
				s.addTrollConnection(t)

			// del a client
			case t := <-s.delCh:
				s.deleteTrollConnection(t)

			// recieve a message from a client troll
			case msg := <-s.messageCh:
				s.recieveMessage(msg)

			case err := <-s.errCh:
				log.Println("Error:", err.Error())

			case <-s.doneCh:
				return
		}
	}
}



