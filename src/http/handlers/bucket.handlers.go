package storage

import (
	"data-storage/src/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeBucket(ctx *gin.Context) {
	bucketName := ctx.Param("name")

	storage.MakeBucket(bucketName)

	ctx.JSON(http.StatusOK, gin.H{})
}

func ListBucketObjects(ctx *gin.Context) {
	bucketName := ctx.Param("name")
	prefix := ctx.Param("prefix")

	objects, err := storage.ListBucketObjects(bucketName, prefix)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"data": objects})
}
