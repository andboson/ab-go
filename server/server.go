package server

import (
	"log"
	 "net/http"
	"time"
	"github.com/gorilla/websocket"
	"ab-go/requests"
	"ab-go/bindata"
	"ab-go/service"
)

const (
// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Buffered channel of outbound messages.
var	Send chan *requests.Result

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn
}

func Init(){
	addr := ":" + service.Args.Port
	Send = make(chan *requests.Result, 3000)
	http.HandleFunc("/", handleAsset("static/index.html"))
	http.HandleFunc("/main.js", handleAsset("static/main.js"))
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		//c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("", message)
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

type ResultMessage struct {
	Ts int64
	Avg string
	Max string
	Min string
	Rps string
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message := <-Send:
			res := &ResultMessage{Ts:time.Now().Unix() * 1000, Avg:message.Avg, Max: message.Max, Min:message.Min, Rps:message.Rps}
			c.ws.WriteJSON(res)
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &connection{ ws: ws}
	go c.writePump()
	c.readPump()
}


//handle static assets
func handleAsset(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := bindata.Asset(path)
		if err != nil {
			log.Fatal(err)
		}

		n, err := w.Write(data)
		if err != nil {
			log.Fatal(err)
		}

		if n != len(data) {
			log.Fatal("wrote less than supposed to")
		}
	}
}