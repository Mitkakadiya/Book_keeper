package main

type Status string

const (
	Read    Status = "read"
	Reading Status = "reading"
	ToRead  Status = "to_read"
)

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" validate:"required" gorm:"index"`
	Status string `json:"status" gorm:"default:to_read"`
	Author string `json:"author" validate:"required" gorm:"index"`
	Year   int    `json:"year" validate:"required" gorm:"index"`
	UserID int    `json:"user_id"`
}

type Result struct {
	// Username string
	// Title    string
	// Year     int
	// UserID     int `json:"user_id"`
	Status     string
	TotalBooks int `json:"total_books"`
}
