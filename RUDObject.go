package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

var minioClient *minio.Client

func makeBucket(c *gin.Context, destination string) {
	if _, err := minioClient.BucketExists(context.Background(), destination); err != nil {
		if errResp := minio.ToErrorResponse(err); errResp.Code == "NoSuchBucket" {
			err = minioClient.MakeBucket(context.Background(), destination, minio.MakeBucketOptions{})
			if err != nil {
				return 
			}
		}
	}
}

func listBucketObjects(c *gin.Context) {
	bucketName := c.Param("name")
	objectsCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
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

func deleteObject(c *gin.Context) {
	bucketName := c.Param("name")
	objectName := c.Param("objectName")
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Object %s deleted successfully", objectName)})
}

func renameObject(c *gin.Context) {
	var renameRequest struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}
	if err := c.BindJSON(&renameRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bucketName := c.Param("name")
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: renameRequest.OldName,
	}
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: renameRequest.NewName,
	}

	_, err := minioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = minioClient.RemoveObject(context.Background(), bucketName, renameRequest.OldName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object renamed successfully"})
}

func copyObjectToBucket(c *gin.Context) {
	bucketName := c.Param("name")
	destination := c.Param("destination")
	objectName := c.Param("object")
	fmt.Println("copie")
	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: objectName,
	}
	dst := minio.CopyDestOptions{
		Bucket: destination,
		Object: objectName,
	}
	fmt.Println("copie",src,dst)
	// if err := makeBucket(c, destination); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"message": "Copying object..."})
	_, err := minioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object copied successfully"})
}


