package mock_data

import (
	"log"
	"time"

	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/infrastructure/database"

	"gorm.io/gorm"
)

func stringPtr(s string) *string {
	return &s
}

type MockDataService struct {
	db *database.DB
}

func NewMockDataService(db *database.DB) *MockDataService {
	return &MockDataService{
		db: db,
	}
}

func (m *MockDataService) WipeAndSeedData() error {
	log.Println("🧹 Wiping existing data...")

	// Clear all data from tables (in correct order due to foreign keys)
	if err := m.db.Exec("DELETE FROM targets").Error; err != nil {
		return err
	}
	if err := m.db.Exec("DELETE FROM missions").Error; err != nil {
		return err
	}
	if err := m.db.Exec("DELETE FROM spy_cats").Error; err != nil {
		return err
	}

	// Reset auto-increment sequences
	if err := m.db.Exec("ALTER SEQUENCE targets_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Could not reset targets sequence: %v", err)
	}
	if err := m.db.Exec("ALTER SEQUENCE missions_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Could not reset missions sequence: %v", err)
	}
	if err := m.db.Exec("ALTER SEQUENCE spy_cats_id_seq RESTART WITH 1").Error; err != nil {
		log.Printf("Warning: Could not reset spy_cats sequence: %v", err)
	}

	log.Println("🌱 Seeding mock data...")

	return m.db.Transaction(func(tx *gorm.DB) error {
		// Create spy cats with valid TheCatAPI breed names
		cats := []entities.SpyCat{
			{
				Name:              "Jane",
				YearsOfExperience: 5,
				Breed:             "Abyssinian",
				Salary:            75000,
			},
			{
				Name:              "Kvas",
				YearsOfExperience: 3,
				Breed:             "Maine Coon",
				Salary:            65000,
			},
			{
				Name:              "Mittens",
				YearsOfExperience: 7,
				Breed:             "Siamese",
				Salary:            85000,
			},
			{
				Name:              "Luna",
				YearsOfExperience: 2,
				Breed:             "Persian",
				Salary:            55000,
			},
			{
				Name:              "Felix",
				YearsOfExperience: 4,
				Breed:             "Bengal",
				Salary:            70000,
			},
		}

		for i := range cats {
			if err := tx.Create(&cats[i]).Error; err != nil {
				return err
			}
		}

		now := time.Now()

		// Create missions in different states
		missions := []entities.Mission{
			// Active mission with cat assigned
			{
				Name:        "Operation Goldfish",
				Description: "Infiltrate the aquarium and gather intelligence on the rare goldfish smuggling operation.",
				StartDate:   now.AddDate(0, 0, -5),
				EndDate:     now.AddDate(0, 0, 10),
				CatID:       &cats[0].ID, // Shadow
				IsCompleted: false,
			},
			// Pending mission (no cat assigned)
			{
				Name:        "Mission Catnip Cartel",
				Description: "Investigate the underground catnip distribution network in the city.",
				StartDate:   now.AddDate(0, 0, 2),
				EndDate:     now.AddDate(0, 0, 20),
				IsCompleted: false,
			},
			// Completed mission
			{
				Name:        "Operation Mouse Hunt",
				Description: "Successfully completed mission to eliminate the mouse infestation in the warehouse district.",
				StartDate:   now.AddDate(0, 0, -30),
				EndDate:     now.AddDate(0, 0, -10),
				IsCompleted: true,
				CompletedAt: func() *time.Time { t := now.AddDate(0, 0, -12); return &t }(),
			},
			// Another active mission with different cat
			{
				Name:        "Project Yarn Ball",
				Description: "Undercover operation to infiltrate the yarn manufacturing facility and uncover quality control secrets.",
				StartDate:   now.AddDate(0, 0, -3),
				EndDate:     now.AddDate(0, 0, 15),
				CatID:       &cats[2].ID, // Mittens
				IsCompleted: false,
			},
		}

		for i := range missions {
			if err := tx.Create(&missions[i]).Error; err != nil {
				return err
			}
		}

		// Update cats with mission assignments
		if err := tx.Model(&cats[0]).Update("mission_id", missions[0].ID).Error; err != nil {
			return err
		}
		if err := tx.Model(&cats[2]).Update("mission_id", missions[3].ID).Error; err != nil {
			return err
		}

		// Create targets for missions
		targets := []entities.Target{
			// Targets for Operation Goldfish (different statuses)
			{
				MissionID: missions[0].ID,
				Name:      "Dr. Fisherman",
				Country:   "Monaco",
				Status:    "completed",
				Notes:     stringPtr("Successfully infiltrated his office and retrieved documents."),
			},
			{
				MissionID: missions[0].ID,
				Name:      "Captain Aquarius",
				Country:   "Greece",
				Status:    "in_progress",
				Notes:     stringPtr("Currently tracking his movements near the harbor."),
			},
			{
				MissionID: missions[0].ID,
				Name:      "Marina Scales",
				Country:   "Italy",
				Status:    "init",
				Notes:     nil,
			},

			// Targets for Mission Catnip Cartel
			{
				MissionID: missions[1].ID,
				Name:      "Pablo Whiskers",
				Country:   "Colombia",
				Status:    "init",
				Notes:     nil,
			},
			{
				MissionID: missions[1].ID,
				Name:      "El Gato",
				Country:   "Mexico",
				Status:    "init",
				Notes:     nil,
			},

			// Targets for completed mission (all completed)
			{
				MissionID: missions[2].ID,
				Name:      "Rodent King",
				Country:   "USA",
				Status:    "completed",
				Notes:     stringPtr("Mission accomplished. Warehouse secured."),
			},

			// Targets for Project Yarn Ball
			{
				MissionID: missions[3].ID,
				Name:      "Ms. Knittington",
				Country:   "UK",
				Status:    "completed",
				Notes:     stringPtr("Obtained yarn quality samples successfully."),
			},
			{
				MissionID: missions[3].ID,
				Name:      "Thread Master",
				Country:   "India",
				Status:    "in_progress",
				Notes:     stringPtr("Infiltrating the textile factory as planned."),
			},
		}

		for i := range targets {
			if err := tx.Create(&targets[i]).Error; err != nil {
				return err
			}
		}

		log.Printf("✅ Successfully seeded:")
		log.Printf("   - %d spy cats", len(cats))
		log.Printf("   - %d missions (1 completed, 2 active with cats, 1 pending)", len(missions))
		log.Printf("   - %d targets", len(targets))

		return nil
	})
}
