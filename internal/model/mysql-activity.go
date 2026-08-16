package model

import "time"

// MySQLTarget describes one MySQL server instance that Devopin monitors.
type MySQLTarget struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Name               string     `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Host               string     `gorm:"size:255;not null" json:"host"`
	Port               int        `gorm:"not null" json:"port"`
	Database           string     `gorm:"size:255;not null" json:"database"`
	Username           string     `gorm:"size:255;not null" json:"username"`
	PasswordCiphertext string     `gorm:"type:text;not null" json:"-"`
	SSLMode            string     `gorm:"size:32;not null" json:"ssl_mode"`
	Enabled            bool       `gorm:"not null" json:"enabled"`
	LastCheckedAt      *time.Time `json:"last_checked_at,omitempty"`
	LastError          string     `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// MySQLActivity is a sampled row from SHOW FULL PROCESSLIST. The MySQL
// connection ID is used as the process identifier.
type MySQLActivity struct {
	ID               string     `gorm:"primaryKey" json:"-"`
	ActivityKey      string     `gorm:"size:255;not null;index" json:"activity_key"`
	TargetID         uint       `gorm:"not null;index" json:"target_id"`
	TargetName       string     `gorm:"size:255;not null;index" json:"target_name"`
	ServerVersion    string     `gorm:"size:128" json:"server_version"`
	PID              int64      `gorm:"not null;index" json:"pid"`
	DatabaseName     string     `gorm:"size:255;index" json:"database_name"`
	Username         string     `gorm:"size:255;index" json:"username"`
	ApplicationName  string     `gorm:"size:255" json:"application_name,omitempty"`
	ClientAddress    string     `gorm:"size:255" json:"client_address,omitempty"`
	Command          string     `gorm:"size:64;index" json:"command"`
	State            string     `gorm:"size:255;index" json:"state"`
	Query            string     `gorm:"type:text" json:"query,omitempty"`
	QueryFingerprint string     `gorm:"size:64;index" json:"query_fingerprint,omitempty"`
	QueryStart       *time.Time `json:"query_start,omitempty"`
	QueryDurationMS  int64      `json:"query_duration_ms"`
	ObservedAt       time.Time  `gorm:"not null;index" json:"observed_at"`
}
