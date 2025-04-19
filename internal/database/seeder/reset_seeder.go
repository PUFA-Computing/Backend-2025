package seeder

import (
	"Backend/internal/database"
	"context"
	"log"
)

func ResetDatabase() error {
	// Order matters due to foreign key constraints
	// Delete from child tables first
	tables := []string{
		"news_likes",
		"aspirations_upvote",
		"event_registrations",
		"news",
		"events", 
		"aspirations",
		"role_permissions",
		"permissions",
		"users",
		"roles",
		"organizations",
	}

	for _, table := range tables {
		log.Printf("Truncating table: %s", table)
		_, err := database.DB.Exec(context.Background(), "TRUNCATE TABLE "+table+" CASCADE")
		if err != nil {
			return err
		}
	}

	// Reset sequences
	sequences := []string{
		"news_id_seq",
		"events_id_seq",
		"aspirations_id_seq",
		"permissions_id_seq",
		"roles_id_seq",
		"organizations_id_seq",
	}

	for _, seq := range sequences {
		_, err := database.DB.Exec(context.Background(), "ALTER SEQUENCE "+seq+" RESTART WITH 1")
		if err != nil {
			return err
		}
	}

	return nil
}
