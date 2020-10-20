package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func main() {
	router := gin.Default()
	router.GET("/race", joinRace)
	router.Run()
}

func checkOrigin(r *http.Request) bool {
	return true
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		var message = []byte("[Server] Connection established.")
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func joinRace(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	reader(conn)
}
