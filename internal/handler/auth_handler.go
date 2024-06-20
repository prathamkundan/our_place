package handler

import (
	"context"
	"io"
	"log"
	"net/http"
	"space/internal/service"

	"github.com/gin-gonic/gin"
)

func HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := service.GoogleOAuthConfig.Exchange(context.Background(), code)
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
