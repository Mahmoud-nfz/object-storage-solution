package handlers

import (
	"data-storage/src/auth"
	"data-storage/src/config"
	"data-storage/src/storage"
	"data-storage/src/utils"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"path"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/minio/minio-go/v7"
)

var chunkFilenameFormat = "chunk_%d"

type ChunkMessage struct {
	Index uint64 `json:"index"`
	Data  string `json:"data"`
}

type InitMessage struct {
	ChunkSize     int64    `json:"chunkSize"`
	MissingChunks []uint64 `json:"missingChunks"`
	NumChunks     uint64   `json:"numChunks"`
}

type CommitFilePayload struct {
	Token string `json:"token"`
}

func WebsocketReceiveObjectHandler(ctx *gin.Context) {
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

	token, ok := ctx.MustGet("token").(string)
	if !ok {
		log.Panicln("Failed to get token from context")
		return
	}

	dataBucket := fmt.Sprintf("data-%s", fileInfo.DataCollectionID)
	if err := storage.MakeBucket(dataBucket); err != nil {
		log.Panicln("Failed to initialize bucket for this file")
		return
	}

	objectPrefix := path.Join(fileInfo.Hash[0:2], fileInfo.Hash)

	numChunks := uint64((fileInfo.Size + config.Env.ChunkSize - 1) / config.Env.ChunkSize)

	// Find missing chunks
	missingChunks := findMissingChunks(objectPrefix, numChunks)

	// Send the initial message to the client
	initMessage := InitMessage{
		ChunkSize:     config.Env.ChunkSize,
		MissingChunks: missingChunks,
		NumChunks:     numChunks,
	}
	if err := conn.WriteJSON(initMessage); err != nil {
		log.Panicln("Error sending initial message: ", err)
		return
	}

	// Begin listening for chunk uploads
	var wg sync.WaitGroup
	wg.Add(len(missingChunks))
	go listenForMessages(conn, objectPrefix, &wg)

	// Wait for all chunks to be uploaded
	log.Println("Waiting for chunks to upload")
	wg.Wait()

	// Combine all chunks
	wg.Add(1)
	log.Println("Done uploading all chunks")

	// Wait for final merge
	go combineAndUploadFile(dataBucket, path.Join(fileInfo.Path, fileInfo.Name), objectPrefix, numChunks, &wg)
	wg.Wait()
	log.Println("Done uploading file")

	// Notify backend of successful file upload
	commitUpload(token)
	log.Println("Done notifying backend")
}

func findMissingChunks(objectPrefix string, numChunks uint64) []uint64 {
	missingChunks := make([]uint64, 0, numChunks)
	existingChunkFiles, err := storage.ListBucketObjects(storage.UploadsBucket, objectPrefix)
	if err != nil {
		log.Panicln("Error listing bucket objects:", err)
		return nil
	}
	existingChunks := make(map[uint64]bool, len(existingChunkFiles))
	for _, chunkFile := range existingChunkFiles {
		var chunkIndex uint64
		format := path.Join(objectPrefix, chunkFilenameFormat)
		_, err := fmt.Sscanf(chunkFile, format, &chunkIndex)
		if err != nil {
			log.Println("Error parsing chunk file name:", err)
			continue
		}
		existingChunks[chunkIndex] = true
	}
	for i := uint64(0); i < numChunks; i++ {
		if _, exists := existingChunks[i]; !exists {
			missingChunks = append(missingChunks, i)
		}
	}
	return missingChunks
}

func listenForMessages(conn *websocket.Conn, objectPrefix string, wg *sync.WaitGroup) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected websocket closure: %v", err)
			}
			break
		}

		go handleMessage(message, objectPrefix, wg)
	}
}

func handleMessage(message []byte, objectPrefix string, wg *sync.WaitGroup) {
	var chunkMsg ChunkMessage
	if err := json.Unmarshal(message, &chunkMsg); err != nil {
		log.Println("JSON unmarshal error: ", err)
		return
	}

	chunkIndex := chunkMsg.Index
	chunkData, err := base64.StdEncoding.DecodeString(chunkMsg.Data)
	if err != nil {
		log.Println("Base64 decode error: ", err)
		return
	}

	log.Println("Handling chunk index ", chunkIndex)

	chunkFilePath := path.Join(objectPrefix, fmt.Sprintf(chunkFilenameFormat, chunkIndex))
	if err := storage.MakeObject(storage.UploadsBucket, chunkFilePath, chunkData); err != nil {
		log.Println("Error writing chunk to file: ", err)
		return
	}

	wg.Done()
}

// combines all chunks from the "uploads" bucket and puts the final file to the "data" bucket.
func combineAndUploadFile(dataBucket, destinationPath, sourcePrefix string, numChunks uint64, wg *sync.WaitGroup) error {
	srcs := make([]minio.CopySrcOptions, numChunks)
	for i := uint64(0); i < numChunks; i++ {
		chunkFilePath := path.Join(sourcePrefix, fmt.Sprintf(chunkFilenameFormat, i))
		srcs[i] = minio.CopySrcOptions{
			Bucket: storage.UploadsBucket,
			Object: chunkFilePath,
		}
	}
	dst := minio.CopyDestOptions{
		Bucket: dataBucket,
		Object: destinationPath,
	}
	_, err := storage.ConcatenateObjects(dst, srcs...)
	if err != nil {
		log.Printf("Failed to compose object: %v", err)
		return err
	}

	log.Printf("Successfully composed and uploaded object: %s/%s", dataBucket, destinationPath)
	wg.Done()
	return nil
}

func commitUpload(jwt string) {
	headers := make(map[string]string)
	headers["X-Api-Key"] = config.Env.APIKey
	headers["Content-Type"] = "application/json"

	payload := CommitFilePayload{
		Token: jwt,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Panicln("Error marshalling to JSON: ", err)
	}

	url, err := url.Parse(config.Env.BackendUrl)
	if err != nil {
		log.Panicln("Error parsing URL: ", err)
	}
	url.Path = "/_files/upload/commit"
	log.Println(url.String())

	statusCode, body, err := utils.SendHTTPRequest("POST", url.String(), headers, data)
	if err != nil {
		log.Panicln("Error commiting upload to backend:", err)
	}
	if statusCode >= 300 {
		log.Panicln("Error commiting upload to backend:", string(body))
	}
}
