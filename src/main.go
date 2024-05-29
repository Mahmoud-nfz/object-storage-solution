package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	router = gin.Default()

	initializeRoutes()

	err := router.Run(":1206")
	if err != nil {
		log.Println("Error starting server:", err)
	}
}
