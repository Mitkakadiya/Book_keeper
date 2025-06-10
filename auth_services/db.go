package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite" // required for driver
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Auto migrate your models
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
	return db
}
