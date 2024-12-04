package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"transaction-service/config"
)

// ConnectDatabase creates and returns a GORM database instance
func ConnectDatabase() *sql.DB {
	config.LoadEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_USER", "postgres"),
		config.GetEnv("DB_PASSWORD", ""),
		config.GetEnv("DB_NAME", "transaction_dev_db"),
		config.GetEnv("DB_PORT", "5432"),
		config.GetEnv("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to the database successfully!")
	return db
}
