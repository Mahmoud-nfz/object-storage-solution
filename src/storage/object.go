package storage

import (
	"context"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func FetchObject(objectName, bucketName string) ([]byte, error) {
	_, err := MinioClient.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	reader, err := MinioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	objectData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return objectData, nil
}

func DeleteObject(c *gin.Context) {
	bucketName := c.Param("name")
	objectName := c.Param("objectName")
	err := MinioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Object deleted successfully"})
}

func RenameObject(c *gin.Context) {
	var renameRequest struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}
	if err := c.BindJSON(&renameRequest); err != nil {
		log.Println("Error:", err)
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
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = MinioClient.RemoveObject(context.Background(), bucketName, renameRequest.OldName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Error:", err)
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
	log.Println("Copying object:", src, dst)
	// if err := makeBucket(c, destination); err != nil {
	// 	log.Println("Error:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	c.JSON(http.StatusOK, gin.H{"message": "Copying object..."})
	_, err := MinioClient.CopyObject(context.Background(), dst, src)
	if err != nil {
		log.Println("Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object copied successfully"})
}
