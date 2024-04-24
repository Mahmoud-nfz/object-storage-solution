package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func MakeBucket(c *gin.Context, destination string) {
	if _, err := MinioClient.BucketExists(context.Background(), destination); err != nil {
		if errResp := minio.ToErrorResponse(err); errResp.Code == "NoSuchBucket" {
			err = MinioClient.MakeBucket(context.Background(), destination, minio.MakeBucketOptions{})
			if err != nil {
				return
			}
		}
	}
}

func ListBucketObjects(c *gin.Context) {
	bucketName := c.Param("name")
	objectsCh := MinioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	var objects []string
	for object := range objectsCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": object.Err.Error()})
			return
		}
		objects = append(objects, object.Key)
	}
	fmt.Println(objects)
	c.JSON(http.StatusOK, gin.H{"objects": objects})
}
