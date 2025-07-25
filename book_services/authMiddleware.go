package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string
		coockieToken := c.Cookies("jwt")
		// take token from coockie if it is not null
		if coockieToken != "" {
			tokenString = coockieToken
		} else {
			//take token from header
			authHeader := c.Get("Authorization")

			if authHeader == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": fiber.StatusUnauthorized,
					"error":  "Unauthorized",
				})
			}

			// split header to get token

			tokenParts := strings.Split(authHeader, " ")

			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": fiber.StatusUnauthorized,
					"error":  "Unauthorized",
				})
			}

			// use token from auth header
			tokenString = tokenParts[1]
		}

		// parse token

		secret := []byte("super-secret-key")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("unexpected sigin method:%v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": fiber.StatusUnauthorized,
				"error":  "Unauthorized",
			})
		}

		userId := token.Claims.(jwt.MapClaims)["userId"]

		Result, err := http.Get(fmt.Sprintf("http://127.0.0.1:3000/auth/users/%v", userId))
		if err != nil {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": fiber.StatusUnauthorized,
				"error":  "Unauthorized",
			})
		}

		var userData map[string]interface{}
		json.NewDecoder(Result.Body).Decode(&userData)
		fmt.Println(userData["id"], userData["username"])

		fmt.Println(userId)
		c.Locals("userId", userId)
		return c.Next()
	}
}
