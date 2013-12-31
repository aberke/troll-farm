package trolls

import (
        "log"
        "net/http"

        "code.google.com/p/go.net/websocket"
)

const NEW_CONNECTION_ENDPOINT = "/connect";


// troll server
type Server struct {
        trolls    		map[int]*Troll
        trollsDataMap	map[int]*TrollData
        addCh     		chan *Troll
        delCh     		chan *Troll
        sendAllCh 		chan *OutgoingMessage
        doneCh    		chan bool
        errCh     		chan error
}

// Create new troll server.
func NewServer() *Server {
	trolls := make(map[int]*Troll)
	trollsDataMap := make(map[int]*TrollData)
	addCh := make(chan *Troll)
	delCh := make(chan *Troll)
	sendAllCh := make(chan *OutgoingMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		trolls,
		trollsDataMap,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}


func (s *Server) AddTrollConnection(t *Troll) {
        s.addCh <- t
}

func (s *Server) Del(t *Troll) {
        s.delCh <- t
}

func (s *Server) SendAll(msg *OutgoingMessage) {
	log.Println("server.SendAll msg: ", msg)
	s.sendAllCh <- msg
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

func (s *Server) recieveTestMessage(trollID int) {
	var msg *OutgoingMessage
	msg = OutgoingTestMessage(trollID)
	s.trolls[trollID].Write(msg)
}
func (s *Server) sendTrollsMessage(trollID int) {
	log.Println("sendTrollsMessage", s.trollsDataMap)
	log.Println("sendTrollsMessage", s.trollsDataMap[1])
	var msg *OutgoingMessage
	msg = OutgoingTrollsMessage(trollID, s.trollsDataMap)
	s.trolls[trollID].Write(msg)
}
func (s *Server) recieveTrollsMessage(trollID int) {
	s.sendTrollsMessage(trollID)
}

func (s *Server) recieveMessageMessage(trollID int, data string) {
	log.Println("TODO: recieveTrollsMessage")
}
func (s *Server) recieveMoveMessage(trollID int, data string) {
	log.Println("TODO: recieveTrollsMessage")
}


// when troll client recieves message, sends the IncomingMessage to server to be handled
func (s *Server) recieveMessage(trollID int, msg *IncomingMessage) {
	log.Println("Server handleMessage from ", trollID)
	log.Println("incoming msg: ", msg)

	switch msg.Type {
	case "test":
		s.recieveTestMessage(trollID)
	case "trolls":
		s.recieveTrollsMessage(trollID)
	case "message":
		s.recieveMessageMessage(trollID, msg.Data)
	case "move":
		s.recieveMoveMessage(trollID, msg.Data)
	default:
		log.Println("Unknown message type recieved: ", msg.Type)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
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
				log.Println("Added new troll")
				s.trolls[t.id] = t
				td := NewTrollData(t)
				s.trollsDataMap[t.id] = td
				log.Println("Now", len(s.trolls), "trolls connected.")
				//s.sendTrollsMessage(t.id)

			// del a client
			case t := <-s.delCh:
				log.Println("Delete troll")
				delete(s.trolls, t.id)

			// broadcast message for all clients
			case msg := <-s.sendAllCh:
				log.Println("Send all:", msg)
				s.sendAll(msg)

			case err := <-s.errCh:
				log.Println("Error:", err.Error())

			case <-s.doneCh:
				return
		}
	}
}



