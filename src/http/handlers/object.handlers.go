package storage

import (
	"data-storage/src/storage"

	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteObject(ctx *gin.Context) {
	bucketName := ctx.Param("name")
	objectName := ctx.Param("objectName")
	err := storage.DeleteObject(bucketName, objectName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Object deleted successfully"})
	}
}

func RenameObject(ctx *gin.Context) {
	var renameRequest struct {
		OldName string `json:"oldName"`
		NewName string `json:"newName"`
	}
	if err := ctx.BindJSON(&renameRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bucketName := ctx.Param("name")

	err := storage.RenameObject(bucketName, renameRequest.OldName, renameRequest.NewName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Object renamed successfully"})
}

func CopyObjectToBucket(c *gin.Context) {
	bucketName := c.Param("name")
	destination := c.Param("destination")
	objectName := c.Param("object")

	err := storage.CopyObjectToBucket(bucketName, destination, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Object copied successfully"})
}
