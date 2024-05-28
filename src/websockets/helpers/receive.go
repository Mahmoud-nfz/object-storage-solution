package websockets

import (
	"data-storage/src/storage"
	"data-storage/src/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleMultipleFilesReception(c *gin.Context, conn *websocket.Conn, tempDir string, combinedFile *os.File) {
	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			log.Println("No more files to recieve:", err)
			break
		}

		err = HandleOneFileReception(c, conn, tempDir, combinedFile, messageType, message)

		if err != nil {
			log.Println("Error handling file reception, continuing to next file:", err)
		}
	}
}

func HandleOneFileReception(c *gin.Context, conn *websocket.Conn, tempDir string, combinedFile *os.File, messageType int, message []byte) error {

	if messageType == websocket.TextMessage {
		received := string(message)
		log.Println("Received:", received)

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
			return err
		}

		_, err = combinedFile.WriteString(fileName + ":\n")
		if err != nil {
			file.Close()
			log.Panicln("Error writing original filename to combined file:", err)
			return err
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte("ready"))
		if err != nil {
			file.Close()
			log.Panicln("Error sending ready message:", err)
			return err
		}

		for {
			messageType, content, err := conn.ReadMessage()

			if err != nil {
				log.Println("Connection closed by client")
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
	return nil
}
