package seeder

import (
	"Backend/internal/database"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

func SeedNews() error {
	// Sample news data
	news := []struct {
		title          string
		content        string
		publishDate    time.Time
		slug           string
		thumbnail      string
		organizationID int
	}{
		{
			title:       "Welcome to PUFA Computing",
			content:     "PUFA Computing is proud to announce our new website launch, providing better service to our members...",
			publishDate: time.Now(),
			slug:        "welcome-to-pufa-computing",
			thumbnail:   "/images/news/cs-program.jpg",
			organizationID: 1,
		},
		{
			title:       "Upcoming Technology Workshop Series",
			content:     "Join us for a series of workshops covering the latest technologies in software development...",
			publishDate: time.Now().Add(-24 * time.Hour),
			slug:        "upcoming-technology-workshop-series",
			thumbnail:   "/images/news/cs-program2.jpg",
			organizationID: 1,
		},
	}

	// Get admin user ID (assuming there's at least one admin user)
	var adminID uuid.UUID
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM users LIMIT 1").Scan(&adminID)
	if err != nil {
		return err
	}

	// Insert news
	for _, n := range news {
		_, err := database.DB.Exec(context.Background(), `
			INSERT INTO news (title, content, user_id, publish_date, thumbnail, slug, organization_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, n.title, n.content, adminID, n.publishDate, n.thumbnail, n.slug, n.organizationID)
		
		if err != nil {
			log.Printf("Error seeding news %s: %v", n.title, err)
			return err
		}
	}

	return nil
}
