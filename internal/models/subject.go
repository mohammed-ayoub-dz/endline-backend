package models

type Subject struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Title string `gorm:"not null" json:"title"`

Videos []Video `gorm:"foreignKey:SubjectID" json:"videos,omitempty"`
}