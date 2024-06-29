package ws

import (
	"time"
	"workwithimages/domain/models"
)

type WServer struct {
	Clients   map[*Client]bool
	Add       chan *Client
	Delete    chan *Client
	Broadcast chan *models.Message
}

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func NewWSevrver() *WServer {
	return &WServer{
		Clients:   make(map[*Client]bool),
		Add:       make(chan *Client),
		Delete:    make(chan *Client),
		Broadcast: make(chan *models.Message),
	}
}

func (ws *WServer) WSRun() {
	for {
		select {
		case msg := <-ws.Broadcast:
			go ws.SendToAll(msg)
		case client := <-ws.Add:
			ws.Clients[client] = true
		case client := <-ws.Delete:
			client.Conn.Close()
			delete(ws.Clients, client)
		}
	}
}

func (ws *WServer) SendToAll(msg *models.Message) {
	for client := range ws.Clients {
		go func() {
			client.Send <- msg
		}()
	}
}
