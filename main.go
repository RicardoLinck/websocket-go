package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func (hub *Hub) serveWs(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go func(conn *websocket.Conn) {
		client := &Client{conn: conn, send: make(chan string)}
		hub.register <- client

		go client.sendData()

	}(conn)
}

func main() {
	hub := newHub()
	go hub.run()
	go hub.generateData()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", hub.serveWs)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
