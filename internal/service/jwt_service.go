package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var InvalidJWTErr = errors.New("Invalid JWT")
var InvalidJWTAlgErr = errors.New("Invalid signing algorithm")
var JWTParseErr = errors.New("Could not parse JWT")

func GenerateJWT(user map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user["id"],
		"name":  user["name"],
		"email": user["email"],
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return tokenString, nil
}

func VerifyJwt(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, InvalidJWTAlgErr
		}
		return []byte(JwtSecret), nil
	})

	if err != nil {
		return nil, JWTParseErr
	}

	// Check token validity
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	} else {
		return nil, InvalidJWTErr
	}
}

func GetNameFromToken(token *jwt.Token) (string, error) {
	// Type assert the claims to jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("unable to parse claims")
	}

	// Extract the name field from claims
	name, ok := claims["name"].(string)
	if !ok {
		return "", fmt.Errorf("name field not found in token claims")
	}

	return name, nil
}
