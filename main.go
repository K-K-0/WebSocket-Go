package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients []*websocket.Conn
var mutex sync.Mutex

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error while connecting", err)
		return
	}

	defer conn.Close()

	mutex.Lock()
	clients = append(clients, conn)
	mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("error while reading messages", err)
		}

		mutex.Lock()

		for _, c := range clients {
			if c != conn {
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("error while writing message", err)

				}
			}
		}
		mutex.Unlock()
	}

	mutex.Lock()
	for i, c := range clients {
		if c == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	mutex.Unlock()
}

func main() {
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})

	r.Run("localhost:8000")
}
