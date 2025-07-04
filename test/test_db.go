package test

import (
	"fmt"
	"log"
	"os"
	"time"

	"api-rbac/db"
)

func TestDB() {
	fmt.Println("Testing database connection...")

	// Set environment
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "local")
	}

	fmt.Printf("GO_ENV: %s\n", os.Getenv("GO_ENV"))

	// Test database connection
	connection := db.Instance()
	if connection == nil {
		log.Fatal("Failed to get database instance")
	}

	// Test the connection
	sqlDB := connection.DB()
	if sqlDB == nil {
		log.Fatal("Failed to get underlying sql.DB")
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping the database
	err := sqlDB.Ping()
	if err != nil {
		log.Printf("Database ping failed: %v", err)
		log.Fatal("Cannot connect to database")
	}

	fmt.Println("Database connection successful!")
	fmt.Println("Connection pool stats:")
	fmt.Printf("  Max Open Connections: %d\n", sqlDB.Stats().MaxOpenConnections)
	fmt.Printf("  Open Connections: %d\n", sqlDB.Stats().OpenConnections)
	fmt.Printf("  In Use: %d\n", sqlDB.Stats().InUse)
	fmt.Printf("  Idle: %d\n", sqlDB.Stats().Idle)
}
