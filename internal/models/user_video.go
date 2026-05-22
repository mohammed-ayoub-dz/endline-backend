package models

import "time"

type UserVideo struct {
    UserID    uint      `gorm:"primaryKey"`
    VideoID   uint      `gorm:"primaryKey"`
    AddedAt   time.Time `gorm:"autoCreateTime"`
    User      User      `gorm:"foreignKey:UserID"`
    Video     Video     `gorm:"foreignKey:VideoID"`
}