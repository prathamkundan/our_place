package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	. "space/internal/core"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var upgrader websocket.Upgrader
var canvas *Canvas
var hub *Hub
var googleOAuthConfig *oauth2.Config

func init() {
	canvas = NewCanvas(512, 512)
	hub = SetupHub()
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment variables from .env")
	}

	googleOAuthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	hub.Register(canvas)
	r := gin.Default()

	r.GET("/ws", func(ctx *gin.Context) {
		handleConnection(ctx.Writer, ctx.Request)
	})
	r.GET("/auth/google", handleGoogleCallback)

	r.Run(":8000")
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection from :", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	_, err = SetupClient(conn, hub, canvas)
	if err != nil {
		log.Printf("Could not set up client for: %s", conn.RemoteAddr())
	}
}

func handleGoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("code exchange failed: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "code exchange failed"})
		return
	}

	userinfoUrl := "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	resp, err := http.Get(userinfoUrl + token.AccessToken)
	if err != nil {
		log.Printf("Could not get user info")
		return
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Printf("%s", userData)
	return
}
