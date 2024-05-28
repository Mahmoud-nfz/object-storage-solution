package main

import (
	"data-storage/src/storage"
	websockets "data-storage/src/websockets/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initializeRoutes() {

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.GET("/download", websockets.WebsocketSendObjectHandler)

	router.GET("/upload", websockets.WebsocketReceiveObjectHandler)

	bucketRoutes := router.Group("/bucket")
	{
		bucketRoutes.GET("/:name/objects", storage.ListBucketObjects)

		bucketRoutes.DELETE("/:name/object/:objectName", storage.DeleteObject)

		bucketRoutes.POST("/:name/object/rename", storage.RenameObject)

		bucketRoutes.POST("/:name/:destination/:object", storage.CopyObjectToBucket)
	}
}
