package ws

import (
	"github.com/gorilla/websocket"
	"time"
	"workwithimages/domain/models"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan *models.Message)
var Delete = make(chan *websocket.Conn)
var Add = make(chan *websocket.Conn)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func WSRun() {
	for {
		select {
		case msg := <-Broadcast:
			go SendToAll(msg)
		case conn := <-Add:
			Clients[conn] = true
		case conn := <-Delete:
			conn.Close()
			delete(Clients, conn)
		}
	}
}

func ReadPump(conn *websocket.Conn) {
	Add <- conn
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(t string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			Delete <- conn
			return
		}
		Broadcast <- &msg
	}
}

func WritePump(conn *websocket.Conn) {
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

func SendToAll(msg *models.Message) {
	for client := range Clients {
		go func(conn *websocket.Conn) {
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := conn.WriteJSON(msg)
			if err != nil {
				conn.Close()
			}
		}(client)
	}
}
