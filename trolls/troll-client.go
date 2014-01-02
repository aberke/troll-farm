package trolls

import (
        "fmt"
        "io"
        "log"
        "os"
        "math"
        "encoding/json"

        "code.google.com/p/go.net/websocket"
)

const channelBufSize = 5

var maxId int = 0



// Troll client JSON data
type TrollData struct {
    Name        string  // used names: {"DELETE": indicates to client to delete troll}
    Color       string
    Coordinates map[string]int
    Messages    []string
    Points      int64
}
// Create new TrollData from Troll
func NewTrollData(troll *Troll) *TrollData {
    log.Println("*** NewTrollData *****")

    coordinates     := make(map[string]int)
    coordinates["x"] = int(math.Mod(float64(troll.id), 9))
    coordinates["y"] = 0
    messages        := make([]string, 5)
    td := TrollData{"no-name", "#FF00FF", coordinates, messages, 0}

    encodedTd, err := json.MarshalIndent(td, "", " ")
    if err != nil {
        fmt.Println("****** err *****", err)
    }
    os.Stdout.Write(encodedTd)

    return &td
}


// Troll client.
type Troll struct {
    id     int
    ws     *websocket.Conn
    server *Server
    ch     chan *OutgoingMessage
    doneCh chan bool
}

// Create new chat client.
func NewTroll(ws *websocket.Conn, server *Server) *Troll {

    if ws == nil {
        panic("ws cannot be nil")
    }
    if server == nil {
        panic("server cannot be nil")
    }

    maxId++
    ch := make(chan *OutgoingMessage, channelBufSize)
    doneCh := make(chan bool)

    return &Troll{maxId, ws, server, ch, doneCh}
}

func (t *Troll) Write(msg *OutgoingMessage) {
    select {
        case t.ch <- msg:
        default:
            t.server.Del(t)
            err := fmt.Errorf("troll client %d is disconnected.", t.id)
            t.server.Err(err)
    }
}

func (t *Troll) Done() {
    t.doneCh <- true
}

// Listen Write and Read request via chanel
func (t *Troll) Listen() {
    go t.listenWrite()
    t.listenRead()
}

// Listen write request via chanel
func (t *Troll) listenWrite() {
    for {
        select {

            // send message to the client
            case msg := <-t.ch:
                websocket.JSON.Send(t.ws, msg)

            // receive done request
            case <-t.doneCh:
                t.server.Del(t)
                t.doneCh <- true // for listenRead method
                return
        }
    }
}

// Listen read request via chanel
func (t *Troll) listenRead() {
    for {
        select {

            // receive done request
            case <-t.doneCh:
                t.server.Del(t)
                t.doneCh <- true // for listenWrite method
                return

            // read data from websocket connection
            default:
                var msg IncomingMessage
                err := websocket.JSON.Receive(t.ws, &msg)
                if err == io.EOF {
                    t.doneCh <- true
                } else if err != nil {
                    t.server.Err(err)
                } else {
                    msg.LocalTroll = t.id
                    t.server.RecieveMessage(&msg)
                }
        }
    }
}











