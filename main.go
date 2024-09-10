package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v\n", err)
	}

	// Fetch the PostgreSQL connection string from environment variable
	connString := os.Getenv("DB_URL")
	if connString == "" {
		log.Fatalf("DB_URL environment variable is not set")
	}

	// Connect to the PostgreSQL database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Connected to the database successfully!")

	// Create table SQL query (if it doesn't already exist)
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			age INT NOT NULL,
			email TEXT UNIQUE NOT NULL
		);
	`

	// Execute the query to create the table
	_, err = conn.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatalf("Failed to execute create table query: %v\n", err)
	}
	fmt.Println("Table created successfully!")

	// Insert value SQL query
	insertSQL := `
		INSERT INTO users (name, age, email) 
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	// Insert a user into the table
	var userID int
	err = conn.QueryRow(context.Background(), insertSQL, "John Doe", 30, "john.doe@example.com").Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to insert user: %v\n", err)
	}

	fmt.Printf("User inserted successfully with ID %d!\n", userID)
}
