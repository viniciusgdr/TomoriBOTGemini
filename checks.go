package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CheckPathDb() {
	pathExists := func(path string) bool {
		_, err := os.Stat(path)
		return !os.IsNotExist(err)
	}
	if !pathExists("db") {
		os.Mkdir("db", 0755)
	}
	if !pathExists("assets/temp") {
		os.Mkdir("assets/temp", 0755)
	}
	if !pathExists("src/commands/maintenance") {
		os.Mkdir("src/commands/maintenance", 0755)
	}
}

func ClearCache() {
	fmt.Println("Cleaning cache...")
	os.RemoveAll("assets/temp")
	os.Mkdir("assets/temp", 0755)
	fmt.Println("Cache cleaned!")
}


func LoadDB() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TIMEZONE"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}