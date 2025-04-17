package seeder

import (
	"Backend/internal/database"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

func SeedEvents() error {
	// Sample events data
	events := []struct {
		title          string
		description    string
		startDate      time.Time
		endDate        time.Time
		userID         uuid.UUID
		status         string
		slug           string
		thumbnail      string
		organizationID int
		maxRegistration int
	}{
		{
			title:       "Compsphere 2025",
			description: "Annual Events from Faculty of Computer Science",
			startDate:   time.Now().AddDate(0, 1, 0),
			endDate:     time.Now().AddDate(0, 1, 2),
			status:      "Open",
			slug:        "compsphere-2024",
			thumbnail:   "/images/workshop.jpg",
			organizationID: 1,
			maxRegistration: 50,
		},
		{
			title:       "CompDay 2025",
			description: "Celebration of Computer Science Birth Day",
			startDate:   time.Now().AddDate(0, 2, 0),
			endDate:     time.Now().AddDate(0, 2, 3),
			status:      "Upcoming",
			slug:        "compday-2024",
			thumbnail:   "/images/conference.jpg",
			organizationID: 1,
			maxRegistration: 100,
		},
	}

	// Get admin user ID (assuming there's at least one admin user)
	var adminID uuid.UUID
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM users LIMIT 1").Scan(&adminID)
	if err != nil {
		return err
	}

	// Insert events
	for _, event := range events {
		event.userID = adminID
		_, err := database.DB.Exec(context.Background(), `
			INSERT INTO events (title, description, start_date, end_date, user_id, status, slug, thumbnail, organization_id, max_registration)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, event.title, event.description, event.startDate, event.endDate, event.userID, event.status, event.slug, event.thumbnail, event.organizationID, event.maxRegistration)
		
		if err != nil {
			log.Printf("Error seeding event %s: %v", event.title, err)
			return err
		}
	}

	return nil
}
