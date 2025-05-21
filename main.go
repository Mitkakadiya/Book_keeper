package main

import "github.com/gofiber/fiber/v2"

func main() {

	//initialize the database
	db := InitializeDB()

	//create fiber app instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	// define the Auth routes. Those will be public
	AuthHandlers(app.Group("/auth"), db)
	// AuthMiddleware(db)
	BookHandler(app.Group("/book", AuthMiddleware(db)), db)

	// verify jwt tokens

	// start server on port 3000
	app.Listen(":3000")
}
