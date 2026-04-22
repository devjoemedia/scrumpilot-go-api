package database

import (
	"fmt"
	"log"

	"github.com/devjoemedia/scrumpilot-go-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	DB = db
	log.Println("✅ Database connected successfully")

	// Creates table for models in db
	if err := Migrate(); err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}
	fmt.Println("✅ migrated database successfully")
}
