package handler

import (
	"log"
	"net/http"
	"space/internal/core"
	"space/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

var resetCookie = &http.Cookie{
	Name:     "token",
	Value:    "",
	Path:     "/",
	Expires:  time.Unix(0, 0),
	MaxAge:   -1,
	HttpOnly: true,
}

func HandleWebSocket(ctx *gin.Context) {
	cookie, err := ctx.Request.Cookie("token")
	authStatus := false
	username := ""

	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("Unauthenticated connection from :", ctx.Request.RemoteAddr)
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting Cookies"})
			return
		}
	} else {
		log.Println("Authenticated connection from :", ctx.Request.RemoteAddr)
		token, err := service.VerifyJwt(cookie.Value)
		if err != nil {
			http.SetCookie(ctx.Writer, resetCookie)
		}

		username, err = service.GetNameFromToken(token)
		if err != nil {
			log.Println("What the hell is even that?")
			http.SetCookie(ctx.Writer, resetCookie)
		}
		authStatus = true

		if len(username) > 64 {
			username = username[:64]
		}
	}

	conn, err := service.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	_, err = core.SetupClient(conn, service.HubInst, service.CanvasInst, authStatus, username)
	if err != nil {
		log.Printf("Could not set up client for: %s", conn.RemoteAddr())
	}
}
