package websockets

import (
	"data-storage/src/websockets"
	helpers "data-storage/src/websockets/helpers"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
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

	helpers.HandleMultipleFilesReception(c, conn, tempDir, combinedFile)
}
