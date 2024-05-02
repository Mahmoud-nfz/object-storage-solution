package websockets

import (
	"data-storage/src/storage"
	"data-storage/src/utils"
	"data-storage/src/websockets"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func WebsocketReceiveObjectHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	conn, err := websockets.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicln("Upgrader error:", err)
		return
	}
	defer conn.Close()

	tempDir, err := os.MkdirTemp("", "uploads")
	if err != nil {
		log.Panicln("Error creating temporary directory:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	combinedFileName := "combined_file.txt"
	combinedFilePath := filepath.Join(tempDir, combinedFileName)
	combinedFile, err := os.Create(combinedFilePath)
	if err != nil {
		log.Panicln("Error creating combined file:", err)
		return
	}
	defer combinedFile.Close()

	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			if err != io.EOF {
				log.Panicln("Error reading next object to recieve:", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			received := string(message)
			// log.Println("Received:", received)
			params := strings.Split(received, "/")
			if len(params) != 3 {
				log.Println("Invalid message format")
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid message format, expected bucket/folder/filename, received: " + received,
				})
			}
			bucketName := strings.Split(received, "/")[0]
			fileName := strings.Split(received, "/")[1]

			FolderName := filepath.Base(received)
			log.Println("Bucket:", bucketName)
			log.Println("Filename:", fileName)
			log.Println("Folder:", FolderName)

			if FolderName == fileName {
				FolderName = "/" + fileName
			} else {
				FolderName = FolderName + "/" + fileName
			}

			ext := filepath.Ext(FolderName)
			newFileName := utils.RandomString(10) + ext
			filePath := filepath.Join(tempDir, newFileName)

			file, err := os.Create(filePath)
			if err != nil {
				log.Panicln("Error creating file:", err)
				continue
			}

			_, err = combinedFile.WriteString(fileName + ":\n")
			if err != nil {
				file.Close()
				log.Panicln("Error writing original filename to combined file:", err)
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, []byte("ready"))
			if err != nil {
				file.Close()
				log.Panicln("Error sending ready message:", err)
				continue
			}

			for {
				messageType, content, err := conn.ReadMessage()
				log.Println("Message type: ", messageType)
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						log.Println("Connection closed by client")
					} else {
						log.Println("Error reading file content:", err)
					}
					file.Close()
					break
				}
				if messageType == websocket.CloseMessage {
					log.Println("Connection closed by client")
					file.Close()
					break
				}
				if messageType != websocket.BinaryMessage {
					log.Println("Invalid message type, expected BinaryMessage")
					file.Close()
					break
				}

				_, err = file.Write(content)
				if err != nil {
					log.Println("Error writing file content:", err)
					file.Close()
					break
				}
			}

			err = storage.UploadToMinioFolder(filePath, FolderName, bucketName)
			if err != nil {
				log.Println("Error uploading file to Minio:", err)
			}

			log.Println("File upload completed")
		}
	}
}
