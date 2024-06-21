package service

import (
	"net/http"
	"os"
	"space/internal/core"

	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

var Upgrader websocket.Upgrader

var HubInst *core.Hub

var CanvasInst *core.Canvas

var JwtSecret string

func InitServices() {
	CanvasInst = core.NewCanvas(512, 512)
	HubInst = core.SetupHub()

	HubInst.Register(CanvasInst)

	Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	GoogleOAuthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	JwtSecret = os.Getenv("JWT_SECRET")
}
