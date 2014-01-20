package trolls

import (
        "log"
        "net/http"
        "strconv"

        "code.google.com/p/go.net/websocket"
)

const NEW_CONNECTION_ENDPOINT = "/connect"


// troll server
type Server struct {
	trolls    		map[int]*Troll
	//gridMap  	 	map[int]*Grid
	gridMap			*GridMap
	trollToGrid  	map[int]int
	addCh     		chan *Troll
	delCh     		chan *Troll
	messageCh 		chan *IncomingMessage
	doneCh    		chan bool
	errCh     		chan error
}

// Create new troll server.
func NewServer() *Server {
	trolls 				:= make(map[int]*Troll)

	/* Server starts off with no Grids
		Grids are added and removed based on trolls connecting/disconnecting */
	//gridMap 	:= make(map[int]*Grid)
	var gridMap *GridMap = NewGridMap()

	trollToGrid 		:= make(map[int]int)

	addCh 				:= make(chan *Troll)
	delCh 				:= make(chan *Troll)
	messageCh 			:= make(chan *IncomingMessage)
	doneCh 				:= make(chan bool)
	errCh 				:= make(chan error)

	return &Server{
		trolls,
		gridMap,
		trollToGrid,
		addCh,
		delCh,
		messageCh,
		doneCh,
		errCh,
	}
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
/* send msg to all the trolls in the given grid */
func (s *Server) sendAll(msg *OutgoingMessage, gridId int) {
	for tId, gId := range s.trollToGrid {
		if (gId == gridId) {
			s.trolls[tId].Write(msg)
		}
	}
}
func (s *Server) sendErrorMessage (tId int) {
	msg := OutgoingErrorMessage()
	s.trolls[tId].Write(msg)
}
func (s *Server) sendUpdateMessage(gId int) {
	grid := s.gridMap.Grid(gId)

	var msg *OutgoingMessage
	msg = OutgoingUpdateMessage(0, grid.UpdateMap())
	s.sendAll(msg, gId)

	// clear out updateMap
	grid.ClearUpdateMap()
}
func (s *Server) sendItemsMessage(trollID int) {
	var msg *OutgoingMessage
	gId := s.trollToGrid[trollID]
	msg = OutgoingItemsMessage(trollID, s.gridMap.Grid(gId).ItemsMap())
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
	moveX, errX := strconv.Atoi(data["x"])
	moveY, errY := strconv.Atoi(data["y"])
	if (errX != nil || errY != nil) {
		s.sendErrorMessage(trollID)
		return
	}

	gId := s.trollToGrid[trollID]

	// TODO: HAVE GridMap move Troll 
	/* get back 2 items: gridId indicates which Grid Troll now lives on.  ValidMove is like err */
	validMove := s.gridMap.Grid(gId).MoveTroll(trollID, moveX, moveY)
	if (!validMove) {
		s.sendItemsMessage(trollID)
		return
	}

	s.sendUpdateMessage(gId)	
}


// when troll client recieves message, sends the IncomingMessage to server to be handled
func (s *Server) recieveMessage(msg *IncomingMessage) {
	//log.Println("incoming msg: ", msg)

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

	tId := t.id
	/* Add Troll to a Grid in Grid Map and get back the id of that Grid */
	gId := s.gridMap.AddTroll(tId)

	s.trollToGrid[tId] = gId
	s.trolls[t.id] = t
	s.sendUpdateMessage(gId)

	log.Println("Added new troll to grid", gId, "- Now", len(s.trolls), "trolls connected.")
	s.sendItemsMessage(tId)
}
func (s *Server) deleteTrollConnection(t *Troll) {
	tId := t.id
	gId := s.trollToGrid[tId]

	s.gridMap.DeleteTroll(gId, tId)
	delete(s.trollToGrid, tId)
	delete(s.trolls, tId)

	// send update message
	s.sendUpdateMessage(gId)
	log.Println("Removed troll from grid", gId, "Now", len(s.trolls), "trolls connected.")
}

// Listen and serve - serves client connection and broadcast request.
func (s *Server) Listen() {
	log.Println("Troll server listening........")

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



