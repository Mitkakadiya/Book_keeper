package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	//initialize the database
	db := InitializeDB()

	//create fiber app instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})
	// AuthMiddleware(db)
	BookHandler(app.Group("/book", AuthMiddleware(db)), db)
	// start server on port 3000
	app.Listen(":5000")

}
