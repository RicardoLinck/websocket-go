package main

import (
	"time"
)

// Hub maintains the set of active clients and the ability to broadcast messages to them
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}

	}
}

func (hub *Hub) generateData() {
	ticker := time.NewTicker(time.Second * 3)
	timeoutTicker := time.NewTicker(time.Second * 15)

	for {
		select {
		case tick := <-ticker.C:
			hub.broadcast <- tick.Format(time.RFC3339)
		case <-timeoutTicker.C:
			for client := range hub.clients {
				hub.unregister <- client
			}
		}

	}
}
