package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}


func main() {
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe(":1206", nil))
}
