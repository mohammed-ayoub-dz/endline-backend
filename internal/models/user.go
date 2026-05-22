package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Username  string         `gorm:"unique;not null" json:"username"`
    Password  string         `gorm:"not null" json:"-"`
    Points    int            `gorm:"default:0" json:"points"`
    Level     string         `gorm:"default:beginner" json:"level"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

    VideoProgresses []UserVideoProgress `gorm:"foreignKey:UserID" json:"-"`
}