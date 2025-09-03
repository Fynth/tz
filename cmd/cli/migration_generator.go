package main

import (
	"tz/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=subscriptions port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err)
	}

	// Создаём таблицу
	err = db.AutoMigrate(&models.Subscription{})
	if err != nil {
		log.Error().Err(err)
	}
}
