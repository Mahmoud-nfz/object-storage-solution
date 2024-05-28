package websockets

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var UploadsDir string
var Upgrader websocket.Upgrader

func init() {

	Upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	dir, err := os.MkdirTemp("", "uploads")
	if err != nil {
		log.Fatalln("Error creating temporary directory:", err)
	} else {
		UploadsDir = dir
	}

}
