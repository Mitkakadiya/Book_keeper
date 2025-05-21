package main

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandler(bookGroup fiber.Router, db *gorm.DB) {
	bookGroup.Get("/:id", func(c *fiber.Ctx) error {
		book := new(Book)
		bookId := c.Params("id")

		if bookId == "" {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Book id is required",
				"status": fiber.StatusBadRequest,
			})
		}

		fmt.Println(bookId)
		if err := db.Where("id = ?", bookId).First(&book).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Book not found",
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		return c.Status(fiber.StatusOK).JSON(book)
	})
	bookGroup.Post("/", func(c *fiber.Ctx) error {
		book := new(Book)
		fmt.Println(c.Locals("userId"))
		book.UserID = int(c.Locals("userId").(float64))

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": fiber.StatusBadRequest,
				"error":  err.Error(),
			})
		}

		if err := validator.New().Struct(book); err != nil {
			validators := err.(validator.ValidationErrors)
			var errorMsg []string
			for _, err := range validators {
				switch err.ActualTag() {
				case "required":
					errorMsg = append(errorMsg, fmt.Sprintf("field %s is required", err.Field()))
				default:
					errorMsg = append(errorMsg, fmt.Sprintf("field %s is invalid", err.Field()))
				}
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": fiber.StatusBadRequest,
				"error":  errorMsg,
			})
		}

		if err := db.Create(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": fiber.StatusBadRequest,
				"error":  err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(book)
	})

	bookGroup.Get("/", func(c *fiber.Ctx) error {
		var books []Book
		if err := db.Find(&books).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Book not found",
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Database error",
			})
		}
		// var user []User
		// db.Table("users").Select("users.username", "books.title", "books.year").Joins(
		// 	"inner join books on books.user_id = users.id").Scan(&Result{})
		var result []Result

		// inner join which return all record which match id
		// db.Table("users").Select("users.username", "books.title", "books.year").Joins(
		// "inner join books on books.user_id = users.id").Scan(&result)

		// left join which return all match data from books table
		// db.Table("books").Select("books.title", "books.year", "users.username").Joins(
		// "left join users on  books.user_id = users.id").Scan(&result)

		// left join which return all match data from users table
		// db.Table("users").Select("books.title", "books.year", "users.username").Joins(
		// "left join books on  books.id = users.id").Scan(&result)

		//cross join returns all recoed from both table with joined data
		// db.Table("users").Select("books.title", "books.year", "users.username").Joins(
		// "cross join books").Scan(&result)

		//cross join returns all recoed from both table where data not avilable than returns to null
		// db.Table("users").Select("books.title", "books.year", "users.username").Joins(
		// "full join books on  books.id = users.id").Scan(&result)
		db.AutoMigrate(&User{})

		db.Table("books").
			Select("status, COUNT(*) AS total_books").
			Group("status").
			Scan(&result)

		fmt.Println(result)
		// fmt.Println(books)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": fiber.StatusOK,
			"books":  result,
		})
	})

	bookGroup.Put("/:id", func(c *fiber.Ctx) error {
		bookId := c.Params("id")
		book := new(Book)

		if bookId == "" {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Book id is required",
				"status": fiber.StatusBadRequest,
			})
		}

		if err := c.BodyParser(&book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}
		if err := validator.New().Struct(book); err != nil {
			validators := err.(validator.ValidationErrors)
			var errorMsg []string
			for _, err := range validators {
				switch err.ActualTag() {
				case "required":
					errorMsg = append(errorMsg, fmt.Sprintf("field %s is required", err.Field()))
				default:
					errorMsg = append(errorMsg, fmt.Sprintf("field %s is invalid", err.Field()))
				}
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": fiber.StatusBadRequest,
				"error":  errorMsg,
			})
		}

		var dbBook Book
		if err := db.First(&dbBook, bookId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Book not found",
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		book.ID = dbBook.ID
		book.UserID = int(c.Locals("userId").(float64))
		if err := db.Save(&book).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update book",
			})
		}
		fmt.Println(bookId)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Book updated successfully",
			"book":    book,
		})
	})

	bookGroup.Patch("/:id", func(c *fiber.Ctx) error {
		bookId := c.Params("id")
		book := new(Book)

		if bookId == "" {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Book id is required",
				"status": fiber.StatusBadRequest,
			})
		}

		if err := c.BodyParser(&book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		var dbBook Book
		if err := db.First(&dbBook, bookId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Book not found",
				})
			}
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Database error",
			})
		}

		book.ID = dbBook.ID

		if book.Title != "" {
			dbBook.Title = book.Title
		}
		if book.Author != "" {
			dbBook.Author = book.Author
		}
		if book.Status != "" {
			dbBook.Status = book.Status
		}
		if book.Year != 0 {
			dbBook.Year = book.Year
		}

		book.UserID = int(c.Locals("userId").(float64))
		if err := db.Save(&dbBook).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update book",
			})
		}

		fmt.Println(bookId)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Book updated successfully",
			"book":    dbBook,
		})
	})

	bookGroup.Delete("/:id", func(c *fiber.Ctx) error {
		bookId := c.Params("id")

		if bookId == "" {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  "Book id is required",
				"status": fiber.StatusBadRequest,
			})
		}

		result := db.Delete(&Book{}, bookId)

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete the book",
			})
		}
		if result.RowsAffected == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No book found with the given id",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "book deleted successfully",
		})
	})

}
