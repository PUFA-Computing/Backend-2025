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
			title: "Welcome to PUFA Computing",
			content: "PUFA Computing is proud to officially launch our new website, designed to deliver a better, more user-friendly experience for students, staff, and visitors. With a modern look, faster navigation, and regularly updated content, this platform will be the central hub for all news, events, and resources related to the Faculty of Computer Science. We invite you to explore the new features, discover what’s new, and stay connected with our community.",
			publishDate: time.Now(),
			slug: "welcome-to-pufa-computing",
			thumbnail: "/images/news/cs-program.jpg",
			organizationID: 1,
		},
		{
			title: "Upcoming Technology Workshop Series",
			content: "Get ready to dive into the future of software development! PUFA Computing is excited to present a comprehensive series of hands-on workshops covering trending technologies such as cloud computing, mobile app development, AI integration, and modern web frameworks. These workshops are designed for students of all skill levels and will be led by industry professionals and experienced faculty. Don't miss this chance to upgrade your skills and stay ahead in the tech world.",
			publishDate: time.Now().Add(-24 * time.Hour),
			slug: "upcoming-technology-workshop-series",
			thumbnail: "/images/news/cs-program2.jpg",
			organizationID: 1,
		},
		{
			title: "Student Achievement: International Hackathon Winners",
			content: "We are thrilled to share that a team of our talented students has won first place at the International Student Hackathon 2024, competing against universities from around the world. Their innovative project impressed the judges with its creativity, technical excellence, and real-world applicability. This victory showcases the high quality of education and collaboration fostered at PUFA Computing. Congratulations to the team on this remarkable achievement!",
			publishDate: time.Now().Add(-48 * time.Hour),
			slug: "student-achievement-hackathon",
			thumbnail: "/images/news/hackathon.jpg",
			organizationID: 1,
		},
		{
			title: "New Research Lab Opening",
			content: "PUFA Computing proudly announces the official opening of our new Artificial Intelligence and Machine Learning Research Lab. This cutting-edge facility is equipped with the latest hardware and software tools to support advanced research in AI, data science, robotics, and more. The lab will serve as a hub for student-led research, faculty innovation, and industry collaboration, fostering a culture of discovery and experimentation.",
			publishDate: time.Now().Add(-72 * time.Hour),
			slug: "new-research-lab-opening",
			thumbnail: "/images/news/research-lab.jpg",
			organizationID: 1,
		},
		{
			title: "Industry Partnership Announcement",
			content: "We are excited to announce a strategic partnership between PUFA Computing and several leading technology companies. This collaboration aims to provide students with real-world internship opportunities, guest lectures from industry experts, and collaborative project experiences. By bridging the gap between academia and the tech industry, we hope to better prepare our students for successful and impactful careers in the ever-evolving digital world.",
			publishDate: time.Now().Add(-96 * time.Hour),
			slug: "industry-partnership-announcement",
			thumbnail: "/images/news/partnership.jpg",
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
