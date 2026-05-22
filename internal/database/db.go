package database

import (
    "fmt"
    "log"
    "github.com/mohammed-ayoub-js/endline-backend/internal/config"
    "github.com/mohammed-ayoub-js/endline-backend/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    err = db.AutoMigrate(&models.User{}, &models.RefreshToken{} , &models.Subject{} , &models.UserVideoProgress{} , &models.Video{} )
    if err != nil {
        return nil, err
    }

    log.Println("Database connected and migrated")
    return db, nil
}