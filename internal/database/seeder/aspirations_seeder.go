package seeder

import (
	"Backend/internal/database"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

func SeedAspirations() error {
	// Sample aspirations data
	aspirations := []struct {
		subject        string
		message        string
		anonymous      bool
		organizationID int
		closed         bool
		adminReply     *string
	}{
		{
			subject:        "Website Improvement Suggestion",
			message:        "Can we add a dark mode feature to the website? It would help reduce eye strain during night time usage.",
			anonymous:      false,
			organizationID: 1,
			closed:         false,
			adminReply:     nil,
		},
		{
			subject:        "New Feature Request",
			message:        "Please consider adding a mobile app version of the platform for better accessibility.",
			anonymous:      true,
			organizationID: 1,
			closed:         true,
			adminReply:     strPtr("Thank you for the suggestion! We'll add this to our roadmap for future development."),
		},
	}

	// Get admin user ID (assuming there's at least one admin user)
	var adminID uuid.UUID
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM users LIMIT 1").Scan(&adminID)
	if err != nil {
		return err
	}

	// Insert aspirations
	for _, aspiration := range aspirations {
		_, err := database.DB.Exec(context.Background(), `
			INSERT INTO aspirations (user_id, subject, message, anonymous, organization_id, closed, admin_reply, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
		`, adminID, aspiration.subject, aspiration.message, aspiration.anonymous, aspiration.organizationID, 
		   aspiration.closed, aspiration.adminReply, time.Now())
		
		if err != nil {
			log.Printf("Error seeding aspiration %s: %v", aspiration.subject, err)
			return err
		}
	}

	return nil
}

// Helper function to get pointer to string
func strPtr(s string) *string {
	return &s
}
