package trolls

import (
        "log"
        "net/http"
        "strconv"

        "code.google.com/p/go.net/websocket"
)

const NEW_CONNECTION_ENDPOINT = "/connect";


// troll server
type Server struct {
        trolls    		map[int]*Troll
        trollsDataMap	map[int]*TrollData
        updateMap		map[int]*TrollData
        addCh     		chan *Troll
        delCh     		chan *Troll
        messageCh 		chan *IncomingMessage
        doneCh    		chan bool
        errCh     		chan error
}

// Create new troll server.
func NewServer() *Server {
	trolls := make(map[int]*Troll)
	trollsDataMap := make(map[int]*TrollData)
	updateMap := make(map[int]*TrollData)
	addCh := make(chan *Troll)
	delCh := make(chan *Troll)
	messageCh := make(chan *IncomingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		trolls,
		trollsDataMap,
		updateMap,
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
func (s *Server) sendAll(msg *OutgoingMessage) {
	for _, t := range s.trolls {
		t.Write(msg)
	}
}
func (s *Server) sendTrollsMessage(trollID int) {
	var msg *OutgoingMessage
	msg = OutgoingTrollsMessage(trollID, s.trollsDataMap)
	s.trolls[trollID].Write(msg)
}
func (s *Server) sendUpdateMessage() {
	var msg *OutgoingMessage
	msg = OutgoingUpdateMessage(0, s.updateMap)
	s.sendAll(msg)

	// clear out updateMap
	s.updateMap = make(map[int]*TrollData)
}

func (s *Server) recievePingMessage(trollID int) {
	var msg *OutgoingMessage
	msg = OutgoingPingMessage(trollID)
	s.trolls[trollID].Write(msg)
}
func (s *Server) recieveTrollsMessage(trollID int) {
	s.sendTrollsMessage(trollID)
}

func (s *Server) recieveMessageMessage(trollID int, data map[string]string) {
	log.Println("TODO: recieveMessageMessage")
}
func (s *Server) recieveMoveMessage(trollID int, data map[string]string) {
	x, _ := strconv.Atoi(data["x"])
	y, _ := strconv.Atoi(data["y"])
	
	s.trollsDataMap[trollID].Coordinates["x"] += x
	s.trollsDataMap[trollID].Coordinates["y"] += y

	s.updateMap[trollID] = s.trollsDataMap[trollID]
	s.sendUpdateMessage()
}


// when troll client recieves message, sends the IncomingMessage to server to be handled
func (s *Server) recieveMessage(msg *IncomingMessage) {
	log.Println("incoming msg: ", msg)

	switch msg.Type {
	case "ping":
		s.recievePingMessage(msg.LocalTroll)
	case "trolls":
		s.recieveTrollsMessage(msg.LocalTroll)
	case "message":
		s.recieveMessageMessage(msg.LocalTroll, msg.Data)
	case "move":
		s.recieveMoveMessage(msg.LocalTroll, msg.Data)
	default:
		log.Println("Unknown message type recieved: ", msg.Type)
	}
}
func (s *Server) addTrollConnection(t *Troll) {
	td := NewTrollData(t)
	log.Println("addTrollConnection *****")
	s.updateMap[t.id] = td
	s.sendUpdateMessage()

	s.trolls[t.id] = t
	s.trollsDataMap[t.id] = td
	log.Println("Added new troll - Now", len(s.trolls), "trolls connected.")
	s.sendTrollsMessage(t.id)
}
func (s *Server) deleteTrollConnection(t *Troll) {
	// set troll to be deleted in updateMap
	td := s.trollsDataMap[t.id]
	td.Name = "DELETE"
	s.updateMap[t.id] = td

	// delete troll
	delete(s.trolls, t.id)
	delete(s.trollsDataMap, t.id)

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



