package ffmpeg

import (
	"data-storage/src/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io/ioutil"
	"net/http"
	"strconv"
)


var MinioClient *minio.Client

func TrimVideo(input []byte, outputFilePath string, startIdx, endIdx string) ([]byte, error) {
    startSeconds, _ := strconv.ParseFloat(startIdx, 64)
    endSeconds, _ := strconv.ParseFloat(endIdx, 64)

    err := ffmpeg.Input(string(input)).Trim(ffmpeg.KwArgs{"start": startSeconds, "end": endSeconds}).
        Output(outputFilePath).OverWriteOutput().Run()
    if err != nil {
        return nil, fmt.Errorf("failed to trim video: %v", err)
    }

    trimmedVideo, err := ioutil.ReadFile(outputFilePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read trimmed video file: %v", err)
    }
	fmt.Println("Trimmed video successfully")
    return trimmedVideo, nil
}

func HandleTrimVideo(c *gin.Context) {
	bucketName := c.Param("bucketName")
	objectName := c.Param("objectName")
	startIdx := c.Param("startIdx")
	endIdx := c.Param("endIdx")

	videoData, err := storage.FetchObject(objectName, bucketName)
	if err != nil {
		http.Error(c.Writer, "Failed to fetch video from MinIO", http.StatusInternalServerError)
		return
	}

	outputFilePath := "trimmed.mp4"
	trimmedVideo, err := TrimVideo(videoData, outputFilePath, startIdx, endIdx)
	if err != nil {
		http.Error(c.Writer, "Failed to trim video", http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "video/mp4", trimmedVideo)
}
