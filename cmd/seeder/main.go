package main

import (
	"Backend/configs"
	"Backend/internal/database"
	"Backend/internal/database/seeder"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load config
	config := configs.LoadConfig()

	// Initialize database
	database.Init(config) // Remove error handling since Init() doesn't return anything

	// Run seeders
	log.Println("Starting data seeding...")

	if err := seeder.SeedEvents(); err != nil {
		log.Printf("Error seeding events: %v", err)
		return
	}
	log.Println("Events seeded successfully")

	if err := seeder.SeedAspirations(); err != nil {
		log.Printf("Error seeding aspirations: %v", err)
		return
	}
	log.Println("Aspirations seeded successfully")

	if err := seeder.SeedNews(); err != nil {
		log.Printf("Error seeding news: %v", err)
		return
	}
	log.Println("News seeded successfully")

	log.Println("Data seeding completed successfully")
}
