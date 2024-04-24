package main

import (
	"data-storage/src/routes"
	"data-storage/src/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	var err error
	storage.MinioClient, err = storage.InitializeMinioClient()
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

	r.POST("/bucket/:name/:destination/:object", storage.CopyObjectToBucket)
	r.GET("/download", func(c *gin.Context){ websocketDownloadHandler(c.Writer, c.Request) })
	r.GET("/upload",func(c *gin.Context){ websocketHandler(c.Writer, c.Request) })
	r.GET("/hello",
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})

	r.GET("/bucket/:name/objects", storage.ListBucketObjects)
	r.DELETE("/bucket/:name/object/:objectName", storage.DeleteObject)
	r.POST("/bucket/:name/object/rename", storage.RenameObject)

	err = r.Run(":1206")
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}

}
