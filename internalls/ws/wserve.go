package ws

import (
	"github.com/gorilla/websocket"
	"time"
	"workwithimages/domain/models"
)

type WServer struct {
	Clients   map[*websocket.Conn]bool
	Broadcast chan *models.Message
	Delete    chan *websocket.Conn
	Add       chan *websocket.Conn
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func NewWSevrver() *WServer {
	return &WServer{
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan *models.Message),
		Delete:    make(chan *websocket.Conn),
		Add:       make(chan *websocket.Conn),
	}
}

func (ws *WServer) WSRun() {
	for {
		select {
		case msg := <-ws.Broadcast:
			go ws.SendToAll(msg)
		case conn := <-ws.Add:
			ws.Clients[conn] = true
		case conn := <-ws.Delete:
			conn.Close()
			delete(ws.Clients, conn)
		}
	}
}

func ReadPump(conn *websocket.Conn, ws *WServer) {
	ws.Add <- conn
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(t string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			ws.Delete <- conn
			return
		}
		ws.Broadcast <- &msg
	}
}

func WritePump(conn *websocket.Conn, ws *WServer) {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ticker.Stop()
				conn.Close()
				return
			}
		}
	}
}

func (ws *WServer) SendToAll(msg *models.Message) {
	for client := range ws.Clients {
		go func(conn *websocket.Conn) {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteJSON(msg)
			if err != nil {
				conn.Close()
			}
		}(client)
	}
}
