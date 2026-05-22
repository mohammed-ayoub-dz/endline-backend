package models

import (
    "time"
    "gorm.io/gorm"
)

type UserVideoProgress struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    UserID     uint      `gorm:"not null;index:idx_user_video,unique" json:"user_id"`
    VideoID    uint      `gorm:"not null;index:idx_user_video,unique" json:"video_id"`
    Progress   int       `gorm:"default:0" json:"progress"` 
    Completed  bool      `gorm:"default:false" json:"completed"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
    PointsEarned int     `gorm:"default:0" json:"points_earned"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`

    User  User  `gorm:"foreignKey:UserID" json:"-"`
    Video Video `gorm:"foreignKey:VideoID" json:"video,omitempty"`
}

func (p *UserVideoProgress) BeforeSave(tx *gorm.DB) error {
    if p.Progress >= 100 && !p.Completed {
        p.Completed = true
        p.CompletedAt = &time.Time{}
        *p.CompletedAt = time.Now()
    }
    return nil
}