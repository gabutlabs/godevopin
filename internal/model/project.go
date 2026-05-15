package model

import "time"

type Project struct {
	ID           uint         `gorm:"primaryKey" json:"id"`
	Name         string       `gorm:"not null" json:"name"`
	PathLog      string       `gorm:"not null" json:"path_log"`
	ProjectType  string       `gorm:"not null" json:"project_type"`
	LogFormat    string       `json:"log_format"`
	LogHistories []LogHistory `gorm:"foreignKey:ProjectID" json:"log_histories,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
