package ffmpeg

import (
    "context"
    "log"
    "net/http"
    "os"
    "path/filepath"
    
    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg "github.com/u2takey/ffmpeg-go"
    "data-storage/src/storage"
)

func TranscodeVideo(inputPath, outputPath, codec string) error {
    err := ffmpeg.Input(inputPath).
        Output(outputPath, ffmpeg.KwArgs{"c:v": codec}).
        OverWriteOutput().
        ErrorToStdOut().
        Run()
    if err != nil {
        log.Println("failed to transcode video: ", err)
        return err
    }
    return nil
}

func HandleTranscodeVideo(c *gin.Context) {
    bucketName := c.Param("bucketName")
    objectName := c.Param("objectName")
    outputObjectName := c.Param("outputObjectName")
    // get codec from outputObjectName extension
    codec := "libx264"
    if filepath.Ext(outputObjectName) == ".webm" {
        codec = "libvpx-vp9"
    }

    inputFilePath := filepath.Join(os.TempDir(), "input-video.mp4")
    outputFilePath := filepath.Join(os.TempDir(), "output-video.mp4")

    err := storage.MinioClient.FGetObject(context.Background(), bucketName, objectName, inputFilePath, minio.GetObjectOptions{})
    if err != nil {
        log.Println("failed to download video: ", err)
        return
    }

    err = TranscodeVideo(inputFilePath, outputFilePath, codec)
    if err != nil {
       log.Println("failed to transcode video: ", err)
        return
    }

    _, err = storage.MinioClient.FPutObject(context.Background(), bucketName, outputObjectName, outputFilePath, minio.PutObjectOptions{})
    if err != nil {
       log.Println("failed to upload video: ", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Video transcoded successfully"})
}
