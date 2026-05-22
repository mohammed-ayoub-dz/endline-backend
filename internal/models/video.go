package models

import "time"

type Video struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"not null" json:"title"`
    Description string    `json:"description"`
    URL         string    `gorm:"unique;not null" json:"url"` 
    Thumbnail   string    `json:"thumbnail"`
    SubjectID   uint      `gorm:"not null;index" json:"subject_id"`
    Subject     Subject   `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}