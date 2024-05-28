package main

import (
	"data-storage/src/storage"

	"log"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	
	if _, err := storage.InitializeMinioClient(); err != nil {
		log.Fatalln("Error initializing Minio client:", err)
	}

	router = gin.Default()

	initializeRoutes()

	err := router.Run(":1206")
	if err != nil {
		log.Fatalln("Error starting server:", err)
	}
}
