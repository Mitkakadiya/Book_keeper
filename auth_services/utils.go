package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenrateAccessToken(user *User) (string, error) {
	secret := []byte("super-secret-key")
	method := jwt.SigningMethodHS256
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 168).Unix(),
	}

	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, err
}
