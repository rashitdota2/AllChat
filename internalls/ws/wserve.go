package ws

import (
	"github.com/gorilla/websocket"
	"workwithimages/domain/models"
)

var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan models.Message)

func WSRun() {
	msg := <-Broadcast
	for client := range Clients {
		err := client.WriteJSON(msg)
		if err != nil {
			client.Close()
			delete(Clients, client)
		}
	}
}
