package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err) // Using Printf for non-fatal error
	}

	databaseURL := os.Getenv("DATABASE_URL")
	log.Printf("DATABASE_URL being used: %s", databaseURL) // Log the URL

	var errOpen error
	db, errOpen = sql.Open("postgres", databaseURL)
	if errOpen != nil {
		log.Fatalf("Failed to open PostgreSQL connection: %v", errOpen)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}
	log.Println("Connected to PostgreSQL")
}

func GetDB() *sql.DB {
	return db
}
