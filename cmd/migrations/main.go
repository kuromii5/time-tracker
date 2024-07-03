package main

import (
	"flag"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Parse command-line flags
	migrateCmd := flag.String("migrate", "", "Specify 'up' or 'down' to run migrations")
	flag.Parse()

	// Get db url
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	// Run migrations
	m, err := migrate.New(
		"file://migrations/",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Perform migration based on command-line flag
	switch *migrateCmd {
	case "up":
		MigrateUp(m)
	case "down":
		MigrateDown(m)
	default:
		log.Fatal("Invalid migrate command. Use '--migrate=up' or '--migrate=down'")
	}

	log.Println("Migrations applied successfully")
}

func MigrateUp(m *migrate.Migrate) {
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}

func MigrateDown(m *migrate.Migrate) {
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
}
