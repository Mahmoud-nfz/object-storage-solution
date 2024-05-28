package ffmpeg

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    ffmpeg "github.com/u2takey/ffmpeg-go"
    "data-storage/src/storage"
)

func hmsToSeconds(hms string) (int, error) {
    parts := strings.Split(hms, ":")
    if len(parts) != 3 {
        return 0, fmt.Errorf("invalid time format")
    }
    hours, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, err
    }
    minutes, err := strconv.Atoi(parts[1])
    if err != nil {
        return 0, err
    }
    seconds, err := strconv.Atoi(parts[2])
    if err != nil {
        return 0, err
    }
    return hours*3600 + minutes*60 + seconds, nil
}

func secondsToHMS(seconds int) string {
    hours := seconds / 3600
    minutes := (seconds % 3600) / 60
    secs := seconds % 60
    return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

func TrimVideo(inputPath, outputPath string, start, duration int) error {
    startHMS := secondsToHMS(start)
    durationHMS := secondsToHMS(duration)

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
    startIdx, err := hmsToSeconds(c.Param("startIdx"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start index"})
        return
    }
    endIdx, err := hmsToSeconds(c.Param("endIdx"))
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
