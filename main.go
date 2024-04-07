package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seed)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randGen.Intn(len(charset))]
	}
	return string(b)
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current working directory:", err)
		return
	}

	var fileName string
	var file *os.File

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("Read:", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			// Assuming the first message is the filename
			fileName = string(message)

			ext := filepath.Ext(fileName)

			newFileName := randomString(10) + ext

			// Create a new file with the received filename
			filePath := filepath.Join(cwd+"/uploads", newFileName)
			file, err = os.Create(filePath)
			if err != nil {
				log.Println("Error creating file:", err)
				break
			}
			defer file.Close()
		} else if messageType == websocket.BinaryMessage {
			// Write the received chunk of data to the file
			if file != nil {
				_, err := file.Write(message)
				if err != nil {
					log.Println("Error writing to file:", err)
					break
				}
			} else {
				log.Println("File is not initialized")
				break
			}
		}
	}

	log.Println("File upload completed")
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
