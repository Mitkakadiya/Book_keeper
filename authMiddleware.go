package main

import (
	"fmt"
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
					"status": fiber.StatusOK,
					"error":  "Unauthorized",
				})
			}

			// split header to get token

			tokenParts := strings.Split(authHeader, "")

			if len(tokenParts) != 2 || tokenParts[0] != "barear" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": fiber.StatusOK,
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
				"status": fiber.StatusOK,
				"error":  "Unauthorized",
			})
		}

		userId := token.Claims.(jwt.MapClaims)["userId"]
		if err := db.Model(User{}).Where("id = ?", userId).Error; err != nil {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": fiber.StatusOK,
				"error":  "Unauthorized",
			})
		}
		fmt.Println(userId)
		c.Locals("userId", userId)
		return c.Next()
	}
}
