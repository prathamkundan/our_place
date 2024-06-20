package handler

import (
	"log"
	"net/http"
	"space/internal/core"
	"space/internal/service"

	"github.com/gin-gonic/gin"
)


func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection from :", r.RemoteAddr)
	conn, err := service.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	_, err = core.SetupClient(conn, service.HubInst, service.CanvasInst)
	if err != nil {
		log.Printf("Could not set up client for: %s", conn.RemoteAddr())
	}
}

func HandleWebSocket(ctx *gin.Context) {
	handleConnection(ctx.Writer, ctx.Request)
}
