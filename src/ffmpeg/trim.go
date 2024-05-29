package ffmpeg

import (
    "context"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "data-storage/src/storage"
    "data-storage/src/websockets"
    "data-storage/src/utils"
    "data-storage/src/websockets/handlers"
    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg "github.com/u2takey/ffmpeg-go"
)

func TrimVideo(inputPath, outputPath string, start, duration int) error {
    startHMS := utils.SecondsToHMS(start)
    durationHMS := utils.SecondsToHMS(duration)

    err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": startHMS}).
        Output(outputPath, ffmpeg.KwArgs{"t": durationHMS}).
        OverWriteOutput().
        ErrorToStdOut().
        Run()

    if err != nil {
        log.Println("failed to trim video: ", err)
        return err
    }
    return nil
}

func HandleTrimVideo(c *gin.Context) {
    bucketName := c.Param("bucketName")
    objectName := c.Param("objectName")
    download := c.Param("download")
    save := c.Param("save")
    startIdx, err := utils.HmsToSeconds(c.Param("startIdx"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start index"})
        return
    }
    endIdx, err := utils.HmsToSeconds(c.Param("endIdx"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end index"})
        return
    }
    duration := endIdx - startIdx
    log.Printf("duration: %d", duration)

    inputFilePath := filepath.Join(os.TempDir(), "input-video.mp4")
    outputFilePath := filepath.Join(os.TempDir(), "output-video.mp4")

    err = storage.MinioClient.FGetObject(context.Background(), bucketName, objectName, inputFilePath, minio.GetObjectOptions{})
    if err != nil {
       log.Println("failed to download video: ", err)
       c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download video"})
       return
    }

    err = TrimVideo(inputFilePath, outputFilePath, startIdx, duration)
    if err != nil {
        log.Println("failed to trim video: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to trim video"})
        return
    }

    if save == "true" {
        _, err = storage.MinioClient.FPutObject(context.Background(), bucketName, "trimmed-"+objectName, outputFilePath, minio.PutObjectOptions{})
        if err != nil {
            log.Println("failed to upload video: ", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload video"})
            return
        }
    }

    if download == "true" {
        WebsocketSendTrimmedVideo(c, outputFilePath, bucketName)
    } else {
        c.JSON(http.StatusOK, gin.H{"message": "Video trimmed successfully"})
    }
}

func WebsocketSendTrimmedVideo(c *gin.Context, filePath, bucketName string) {
    w := c.Writer
    r := c.Request
    conn, err := websockets.Upgrader.Upgrade(w, r, nil)

    fileName := filepath.Base(filePath)
    err = handlers.DownloadAndSendFileChunks(conn, fileName, bucketName)
    if err != nil {
        log.Println("Error downloading and sending file chunks:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending file chunks"})
    }
}
