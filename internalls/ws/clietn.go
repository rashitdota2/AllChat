package ws

import (
	"github.com/gorilla/websocket"
	"time"
	"workwithimages/domain/models"
)

type Client struct {
	WsHub *WServer
	Conn  *websocket.Conn
	Send  chan *models.Message
}

func ReadPump(client *Client) {
	client.WsHub.Add <- client
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(t string) error { client.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg models.Message
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			client.WsHub.Delete <- client
			close(client.Send)
			return
		}
		client.WsHub.Broadcast <- &msg
	}
}

func WritePump(client *Client) {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ticker.Stop()
				client.Conn.Close()
				return
			}
		case msg := <-client.Send:
			go func() {
				client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					ticker.Stop()
					client.Conn.Close()
					return
				}
			}()
		}
	}
}
