package handler

import (
	"context"
	"encoding/json"
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
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var user map[string]interface{}
	err = json.Unmarshal(userData, &user)
	if err != nil {
		log.Println("Could not parse google response: ", err)
	}

	jwt, err := service.GenerateJWT(user)
	if err != nil {
		log.Println("Could not generate a JWT: ", err)
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    jwt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	c.Redirect(http.StatusFound, "http://localhost:5173/")
}
