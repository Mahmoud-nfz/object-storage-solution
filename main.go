package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func main() {
	
	var err error
	minioClient, err = initializeMinioClient()
	if err != nil {
		log.Fatalln("Error initializing Minio client:", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	r.POST("/bucket/:name/:destination/:object", copyObjectToBucket)
	r.GET("ws",func(c *gin.Context){ websocketHandler(c.Writer, c.Request) })
	r.GET("/bucket/:name/objects", listBucketObjects)
	r.DELETE("/bucket/:name/object/:objectName", deleteObject)
	r.POST("/bucket/:name/object/rename", renameObject)
	
	err = r.Run(":1206")
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
	
}
