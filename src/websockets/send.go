package websockets

import (
	"bytes"
	"context"
	"data-storage/src/storage"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/minio/minio-go/v7"
)

type Message struct {
	BucketName string `json:"bucketName"`
	FileName   string `json:"fileName"`
}

func WebsocketSendObjectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("Read:", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			var msg Message
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("Error decoding JSON:", err)
				continue
			}

			log.Printf("Download request for Bucket: %s, File: %s\n", msg.BucketName, msg.FileName)

			err = DownloadAndSendFileChunks(conn, msg.FileName, msg.BucketName)
			if err != nil {
				log.Println("Error downloading and sending file chunks:", err)
				continue
			}

			log.Println("File download completed")
			break
		}
	}
}

func DownloadAndSendFileChunks(conn *websocket.Conn, fileName, bucketName string) error {

	var fileChunks [][]byte
	log.Println("Downloading file chunks")
	object, err := storage.MinioClient.GetObject(context.Background(), bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		return err
	}

	buffer := make([]byte, stat.Size)
	n, err := io.ReadFull(object, buffer)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return err
		}
	}
	fileChunks = append(fileChunks, buffer[:n])
	log.Printf("Downloaded %d bytes\n", n)

	combinedFile := bytes.Join(fileChunks, []byte{})
	err = conn.WriteMessage(websocket.BinaryMessage, combinedFile)
	if err != nil {
		return err
	}

	return nil
}