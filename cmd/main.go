package main

import (
	"log"
	"space/internal/handler"
	"space/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment variables from .env")
	}

    service.InitServices()
}

func main() {
	r := gin.Default()

	r.GET("/ws", handler.HandleWebSocket)
	r.GET("/auth/google", handler.HandleGoogleCallback)

	r.Run(":8000")
}


