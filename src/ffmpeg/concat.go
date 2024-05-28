package ffmpeg

import (
    "context"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg "github.com/u2takey/ffmpeg-go"
    "data-storage/src/storage"
)

func ConcatVideos(inputPaths []string, outputPath string) error {
    inputs := make([]*ffmpeg.Stream, len(inputPaths))
    for i, path := range inputPaths {
        inputs[i] = ffmpeg.Input(path).Get("0")
    }

    err := ffmpeg.Concat(inputs).
        Output(outputPath).
        OverWriteOutput().
        ErrorToStdOut().
        Run()
    if err != nil {
        log.Println("failed to concatenate videos: ", err)
        return err
    }
    return nil
}

func HandleConcatVideos(c *gin.Context) {
    bucketName := c.Param("bucketName")
    outputObjectName := c.Param("outputObjectName")

    inputObjectNames := c.QueryArray("inputObjectNames")
    if len(inputObjectNames) < 2 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "At least two input videos are required"})
        return
    }

    tempDir := os.TempDir()
    inputFilePaths := make([]string, len(inputObjectNames))
    for i, objectName := range inputObjectNames {
        inputFilePaths[i] = filepath.Join(tempDir, fmt.Sprintf("input-video-%d.mp4", i))
        err := storage.MinioClient.FGetObject(context.Background(), bucketName, objectName, inputFilePaths[i], minio.GetObjectOptions{})
        if err != nil {
            log.Println("failed to download video ", objectName, ": ", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download video"})
            return
        }
    }

    outputFilePath := filepath.Join(tempDir, "output-video.mp4")

    err := ConcatVideos(inputFilePaths, outputFilePath)
    if err != nil {
        log.Println("failed to concatenate videos: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to concatenate videos"})
        return
    }

    _, err = storage.MinioClient.FPutObject(context.Background(), bucketName, outputObjectName, outputFilePath, minio.PutObjectOptions{})
    if err != nil {
        log.Println("failed to upload video: ", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Videos concatenated successfully"})
}
