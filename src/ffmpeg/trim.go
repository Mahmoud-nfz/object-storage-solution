package ffmpeg

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg "github.com/u2takey/ffmpeg-go"
    "data-storage/src/storage"
    "data-storage/src/utils"
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
        return fmt.Errorf("failed to trim video: %v", err)
    }
    return nil
}

func HandleTrimVideo(c *gin.Context) {
    bucketName := c.Param("bucketName")
    objectName := c.Param("objectName")
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
    fmt.Sprintf("duration: %d", duration)

    inputFilePath := filepath.Join(os.TempDir(), "input-video.mp4")
    outputFilePath := filepath.Join(os.TempDir(), "output-video.mp4")

    err = storage.MinioClient.FGetObject(context.Background(), bucketName, objectName, inputFilePath, minio.GetObjectOptions{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to download video: %v", err)})
        return
    }

    err = TrimVideo(inputFilePath, outputFilePath, startIdx, duration)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to trim video: %v", err)})
        return
    }

    _, err = storage.MinioClient.FPutObject(context.Background(), bucketName, "trimmed-"+objectName, outputFilePath, minio.PutObjectOptions{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload video: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Video trimmed successfully"})
}
