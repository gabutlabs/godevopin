package model

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Path        string    `gorm:"not null" json:"path"`
	ProjectType string    `gorm:"not null" json:"project_type"`
	LogFormat   string    `json:"log_format"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
