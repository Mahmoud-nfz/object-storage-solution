package ffmpeg

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg"github.com/u2takey/ffmpeg-go"
    "data-storage/src/storage"
)

func TrimVideo(inputPath, outputPath string, start, duration int) error {
    err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": start}).
        Output(outputPath, ffmpeg.KwArgs{"t": duration}).
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
    startIdx, err := strconv.Atoi(c.Param("startIdx"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start index"})
        return
    }
    endIdx, err := strconv.Atoi(c.Param("endIdx"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end index"})
        return
    }
    duration := endIdx - startIdx

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