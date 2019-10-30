package main

import "github.com/gorilla/websocket"

// Client represents a client connected via websocket
type Client struct {
	conn *websocket.Conn
	send chan string
}

func (client *Client) sendData() {
	defer client.conn.Close()
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(message))
			// w.Close()
		}
	}
}
