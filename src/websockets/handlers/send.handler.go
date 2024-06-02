package handlers

import (
	"data-storage/src/auth"
	"data-storage/src/storage"
	"fmt"
	"path"

	"context"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/minio/minio-go/v7"
)

func WebsocketSendObjectHandler(ctx *gin.Context) {
	conn, ok := ctx.MustGet("conn").(*websocket.Conn)
	if !ok {
		log.Panicln("Failed to get WebSocket connection from context")
		return
	}
	defer conn.Close()

	fileInfo, ok := ctx.MustGet("claims").(*auth.JWTPayload)
	if !ok {
		log.Panicln("Failed to get claims from context")
		return
	}

	// retrieve file from MinIO
	bucketName := fmt.Sprintf("data-%s", fileInfo.DataCollectionID)
	objectName := path.Join(fileInfo.Path, fileInfo.Name)
	if err := downloadAndSendFile(conn, objectName, bucketName); err != nil {
		log.Panicln("Could not download and send file: ", err)
	}

	log.Println("Done sending file")
}

func downloadAndSendFile(conn *websocket.Conn, objectName, bucketName string) error {
	object, err := storage.MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()

	// Buffer size for each chunk
	const bufferSize = 64 * 1024 // 64KB

	buffer := make([]byte, bufferSize)
	for {
		n, err := object.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
			return err
		}
	}

	return nil
}