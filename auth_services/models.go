package main

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-"`
	Email    string `json:"email" gorm:"unique"` // new column
}
