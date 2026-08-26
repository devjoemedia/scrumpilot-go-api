package database

import "github.com/devjoemedia/scrumpilot-go-api/models"

func Migrate() error {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Todo{},
		&models.Ticket{},
		&models.Comment{},
		&models.RefreshToken{},
	); err != nil {
		return err
	}

	return nil
}
