package main

import (
	"data-storage/src/auth"
	httpHandlers "data-storage/src/http/handlers"
	websocketHandlers "data-storage/src/websockets/handlers"

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

	websocketRoutes := router.Group("/", auth.EnsureUserAuthenticated())
	{
		websocketRoutes.GET("/download", websocketHandlers.WebsocketSendObjectHandler)

		websocketRoutes.GET("/upload", websocketHandlers.WebsocketReceiveObjectHandler)
	}

	bucketRoutes := router.Group("/bucket", auth.EnsureBackendAuthenticated())
	{
		bucketRoutes.GET("/:name/objects/:prefix", httpHandlers.ListBucketObjects)

		bucketRoutes.DELETE("/:name/object/:objectName", httpHandlers.DeleteObject)

		bucketRoutes.POST("/:name/object/rename", httpHandlers.RenameObject)

		bucketRoutes.POST("/:name/:destination/:object", httpHandlers.CopyObjectToBucket)
	}
}
