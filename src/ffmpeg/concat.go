package ffmpeg

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
	
    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg"github.com/u2takey/ffmpeg-go"
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
        return fmt.Errorf("failed to concatenate videos: %v", err)
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
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to download video %s: %v", objectName, err)})
            return
        }
    }

    outputFilePath := filepath.Join(tempDir, "output-video.mp4")

    err := ConcatVideos(inputFilePaths, outputFilePath)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to concatenate videos: %v", err)})
        return
    }

    _, err = storage.MinioClient.FPutObject(context.Background(), bucketName, outputObjectName, outputFilePath, minio.PutObjectOptions{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to upload video: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Videos concatenated successfully"})
}
