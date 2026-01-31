package database

import (
	"log"
	"os"

	"github.com/ahmadeko2017/backed-golang-tugas-1/internal/entity"
	"github.com/ahmadeko2017/backed-golang-tugas-1/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Read database URL from config (Supabase / Postgres)
	dsn := config.GetString("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set. Set it to your Supabase Postgres connection string (e.g. postgres://user:pass@host:5432/dbname).")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database (Postgres): ", err)
	}

	// Determine whether migration is needed:
	// - If tables don't exist -> migrate
	// - If both tables exist but are empty -> migrate (initial data)
	// - Otherwise, skip migration to avoid interfering with existing data
	migrator := DB.Migrator()
	hasCategory := migrator.HasTable(&entity.Category{})
	hasProduct := migrator.HasTable(&entity.Product{})
	needsMigration := false

	if !hasCategory || !hasProduct {
		needsMigration = true
	} else {
		// Both tables exist; check if they're empty
		var catCount int64
		var prodCount int64
		if err := DB.Model(&entity.Category{}).Count(&catCount).Error; err != nil {
			// If counting fails, assume migration is needed
			needsMigration = true
		}
		if err := DB.Model(&entity.Product{}).Count(&prodCount).Error; err != nil {
			needsMigration = true
		}
		if catCount == 0 && prodCount == 0 {
			needsMigration = true
		}
	}

	if needsMigration {
		// Auto Migrate
		err = DB.AutoMigrate(&entity.Category{}, &entity.Product{})
		if err != nil {
			log.Fatal("Failed to migrate database: ", err)
		}
		log.Println("Database migrated successfully (Postgres)")

		// Seed sample data (optional, controlled by configuration / environment variable)
		if config.GetBool("SEED_DATA") {
			SeedData()
			if config.GetBool("SEED_EXIT") {
				log.Println("Seeding completed. Exiting as SEED_EXIT is set to true.")
				os.Exit(0)
			}
		} else {
			log.Println("Database appears empty or uninitialized. To insert sample data, set SEED_DATA=true and restart the application.")
		}
	} else {
		log.Println("Database appears to have existing data; skipping migration and seed to avoid modifying current data")
	}
}
