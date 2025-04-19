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
			title: "Compsphere 2025",
			description: "Compsphere 2025 adalah acara tahunan terbesar dari Fakultas Ilmu Komputer yang menampilkan berbagai kompetisi, seminar teknologi, pameran karya mahasiswa, dan talkshow bersama tokoh-tokoh industri. Acara ini bertujuan untuk menjadi wadah inspirasi dan kolaborasi antar mahasiswa, dosen, dan profesional di bidang teknologi. Jangan lewatkan momen berharga ini untuk belajar, berkembang, dan memperluas relasi.",
			startDate: time.Now().AddDate(0, 1, 0),
			endDate: time.Now().AddDate(0, 1, 2),
			status: "Open",
			slug: "compsphere-2025",
			thumbnail: "/images/events/compsphere.jpg",
			organizationID: 1,
			maxRegistration: 50,
		},
		{
			title: "CompDay 2025",
			description: "CompDay 2025 merupakan perayaan ulang tahun Fakultas Ilmu Komputer yang diselenggarakan setiap tahun dengan semangat kebersamaan dan apresiasi terhadap pencapaian civitas akademika. Acara ini meliputi kegiatan seru seperti perlombaan antar angkatan, awarding night, serta berbagai penampilan dan hiburan. Hadirilah acara penuh keseruan ini untuk merayakan perjalanan dan pencapaian bersama.",
			startDate: time.Now().AddDate(0, 2, 0),
			endDate: time.Now().AddDate(0, 2, 3),
			status: "Upcoming",
			slug: "compday-2024",
			thumbnail: "/images/events/compday.jpg",
			organizationID: 1,
			maxRegistration: 100,
		},
		{
			title: "CSGO Tournament 2025",
			description: "CSGO Tournament 2025 adalah bagian dari rangkaian Sports and Games Olympiad tahunan yang diperuntukkan bagi seluruh mahasiswa Fakultas Ilmu Komputer. Diselenggarakan oleh PUFA Computer Science, turnamen ini menghadirkan kompetisi seru antar tim dalam game Counter-Strike: Global Offensive, menguji strategi, kerjasama tim, dan kemampuan individu. Raih kesempatan menjadi juara dan tunjukkan kemampuanmu di kancah e-sports kampus!",
			startDate: time.Now().AddDate(0, 3, 0),
			endDate: time.Now().AddDate(0, 3, 1),
			status: "Upcoming",
			slug: "tech-career-fair-2024", // ini slug-nya kayaknya typo
			thumbnail: "/images/events/career-fair.jpg",
			organizationID: 1,
			maxRegistration: 200,
		},
		{
			title: "Pre-Bootcamp High School Event",
			description: "Pre-Bootcamp High School Event adalah program intensif selama tiga hari yang dirancang khusus untuk siswa SMA yang tertarik mengenal dunia pemrograman dan teknologi. Dalam program ini, peserta akan belajar dasar-dasar coding, membuat proyek sederhana, dan mendapatkan pembekalan mengenai peluang karir di bidang IT. Kegiatan ini dikemas secara interaktif, menyenangkan, dan sangat cocok untuk pemula yang ingin mulai belajar sejak dini.",
			startDate: time.Now().AddDate(0, 1, 15),
			endDate: time.Now().AddDate(0, 1, 18),
			status: "Upcoming",
			slug: "coding-bootcamp",
			thumbnail: "/images/events/bootcamp.jpg",
			organizationID: 1,
			maxRegistration: 30,
		},
		{
			title: "AI Research Symposium",
			description: "AI Research Symposium adalah simposium tahunan yang menghadirkan presentasi dari peneliti, mahasiswa, dan praktisi dalam bidang kecerdasan buatan (AI). Acara ini menjadi ajang diskusi dan berbagi wawasan seputar inovasi dan penelitian terbaru dalam machine learning, computer vision, NLP, dan topik-topik AI lainnya. Cocok untuk kamu yang tertarik dengan riset teknologi dan ingin membangun koneksi dengan komunitas ilmiah.",
			startDate: time.Now().AddDate(0, 4, 0),
			endDate: time.Now().AddDate(0, 4, 1),
			status: "Upcoming",
			slug: "ai-research-symposium",
			thumbnail: "/images/events/symposium.jpg",
			organizationID: 1,
			maxRegistration: 150,
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
