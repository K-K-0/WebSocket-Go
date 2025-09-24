package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var client = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrader error: ", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	client[conn] = true
	mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read Error:", err)
			break
		}
		fmt.Printf("Received: %s\n", msg)

		mu.Lock()

		for c := range client {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				fmt.Print("Error while reading message")
				c.Close()
				delete(client, c)
			}
		}
		mu.Unlock()
	}
	// mu.Lock()
	// delete(client, conn)
	// mu.Unlock()
}

func main() {
	http.HandleFunc("/ws", HandleWS)

	http.ListenAndServe(":8080", nil)
}
