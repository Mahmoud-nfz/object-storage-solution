package storage

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func DeleteObject(c *gin.Context) {
	bucketName := c.Param("name")
	objectName := c.Param("objectName")
	err := MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Object %s deleted successfully", objectName)})
}

func RenameObject(c *gin.Context) {
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

	_, err := MinioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = MinioClient.RemoveObject(context.Background(), bucketName, renameRequest.OldName, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object renamed successfully"})
}

func CopyObjectToBucket(c *gin.Context) {
	bucketName := c.Param("name")
	destination := c.Param("destination")
	objectName := c.Param("object")

	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: objectName,
	}
	dst := minio.CopyDestOptions{
		Bucket: destination,
		Object: objectName,
	}
	fmt.Println("copy", src, dst)

	c.JSON(http.StatusOK, gin.H{"message": "Copying object..."})
	_, err := MinioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object copied successfully"})
}
