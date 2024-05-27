package main

import (
	"data-storage/src/storage"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var router *gin.Engine

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}

	if _, err := storage.InitializeMinioClient(); err != nil {
		log.Fatalln("Error initializing Minio client:", err)
	}

	router = gin.Default()

	initializeRoutes()

	router.Run(":1206")

}
