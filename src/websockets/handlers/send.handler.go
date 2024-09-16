package handlers

import (
	"data-storage/src/auth"
	"data-storage/src/storage"

	"io"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

	log.Println("Downloading file")
	object, err := storage.GetObject(storage.DataBucket, path.Join(fileInfo.Path, fileInfo.Name))
	if err != nil {
		log.Panicln("Failed to download file", err)
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		log.Panicln("Failed to get file info", err)
		return
	}

	buffer := make([]byte, stat.Size)
	_, err = io.ReadFull(object, buffer)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			log.Panicln("Failed to read the file contents", err)
			return
		}
	}

	log.Println("Sending file")
	err = conn.WriteMessage(websocket.BinaryMessage, buffer)
	if err != nil {
		log.Panicln("Failed to send the file", err)
		return
	}

	log.Println("File download completed")
}
