package websockets

import (
	"data-storage/src/storage"
	"data-storage/src/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"io"
	"github.com/gorilla/websocket"
)


func WebsocketRecieveObjectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()
	tempDir, err := ioutil.TempDir("", "uploads")
	if err != nil {
		log.Println("Error creating temporary directory:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	combinedFileName := "combined_file.txt"
	combinedFilePath := filepath.Join(tempDir, combinedFileName)
	combinedFile, err := os.Create(combinedFilePath)
	if err != nil {
		log.Println("Error creating combined file:", err)
		return
	}
	defer combinedFile.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("Read:", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			received := string(message)
			log.Println("Received:", received)
			bucketName := strings.Split(received, "/")[0]
			fileName := strings.Split(received, "/")[1]

			tempDir, err := ioutil.TempDir("", "uploads")
			if err != nil {
				log.Println("Error creating temporary directory:", err)
				continue
			}
			defer os.RemoveAll(tempDir)

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
				log.Println("Error creating file:", err)
				continue
			}

			_, err = combinedFile.WriteString(fileName + ":\n")
			if err != nil {
				log.Println("Error writing original filename to combined file:", err)
				file.Close()
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, []byte("ready"))
			if err != nil {
				log.Println("Error sending ready message:", err)
				file.Close()
				continue
			}

			for {
				messageType, content, err := conn.ReadMessage()
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
