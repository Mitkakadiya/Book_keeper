package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthHandlers(authGroup fiber.Router, db *gorm.DB) {
	// Define your authentication routes here
	authGroup.Post("/register", func(c *fiber.Ctx) error {
		user := User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username and password required",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		user.Password = string(hashed)
		if err := db.Create(&user).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": fiber.StatusBadRequest,
					"error":  fmt.Sprintf("User already exists with this username: %s", user.Username),
				})
			}
		}

		accessToken, err := GenrateAccessToken(&user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    accessToken,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"accessToken": accessToken,
		})
	})

	authGroup.Post("/login", func(c *fiber.Ctx) error {
		dbUser := new(User)
		authUser := User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if authUser.Username == "" || authUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Username and password required",
			})
		}

		db.Where("username = ?", authUser.Username).First(&dbUser)

		if dbUser.ID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(authUser.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}
		accessToken, err := GenrateAccessToken(dbUser)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    accessToken,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      fiber.StatusOK,
			"accessToken": accessToken,
		})

	})
}
