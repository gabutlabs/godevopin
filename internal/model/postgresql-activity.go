package model

import "time"

// PostgreSQLTarget describes one PostgreSQL cluster that Devopin monitors.
// PasswordCiphertext is intentionally never serialized to API responses.
type PostgreSQLTarget struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Name               string     `gorm:"size:255;not null;uniqueIndex" json:"name"`
	Host               string     `gorm:"size:255;not null" json:"host"`
	Port               int        `gorm:"not null;default:5432" json:"port"`
	Database           string     `gorm:"size:255;not null;default:postgres" json:"database"`
	Username           string     `gorm:"size:255;not null" json:"username"`
	PasswordCiphertext string     `gorm:"type:text;not null" json:"-"`
	SSLMode            string     `gorm:"size:32;not null;default:prefer" json:"ssl_mode"`
	Enabled            bool       `gorm:"not null" json:"enabled"`
	LastCheckedAt      *time.Time `json:"last_checked_at,omitempty"`
	LastError          string     `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// PostgreSQLActivity is a sampled row from pg_stat_activity. The activity
// key combines the target, PID, and backend start time so PID reuse does not
// merge two different sessions in history.
type PostgreSQLActivity struct {
	ID               string     `gorm:"primaryKey" json:"-"`
	ActivityKey      string     `gorm:"size:255;not null;index" json:"activity_key"`
	TargetID         uint       `gorm:"not null;index" json:"target_id"`
	TargetName       string     `gorm:"size:255;not null;index" json:"target_name"`
	ServerVersion    string     `gorm:"size:64" json:"server_version"`
	PID              int64      `gorm:"not null;index" json:"pid"`
	DatabaseName     string     `gorm:"size:255;index" json:"database_name"`
	Username         string     `gorm:"size:255;index" json:"username"`
	ApplicationName  string     `gorm:"size:255" json:"application_name,omitempty"`
	ClientAddress    string     `gorm:"size:255" json:"client_address,omitempty"`
	ClientPort       int        `json:"client_port,omitempty"`
	BackendType      string     `gorm:"size:64" json:"backend_type,omitempty"`
	BackendStart     *time.Time `json:"backend_start,omitempty"`
	TransactionStart *time.Time `json:"transaction_start,omitempty"`
	QueryStart       *time.Time `json:"query_start,omitempty"`
	StateChange      *time.Time `json:"state_change,omitempty"`
	WaitEventType    string     `gorm:"size:64;index" json:"wait_event_type,omitempty"`
	WaitEvent        string     `gorm:"size:128" json:"wait_event,omitempty"`
	BlockingPID      int64      `json:"blocking_pid,omitempty"`
	BlockingDatabase string     `gorm:"size:255;index" json:"blocking_database,omitempty"`
	BlockingUsername string     `gorm:"size:255;index" json:"blocking_username,omitempty"`
	BlockingState    string     `gorm:"size:64" json:"blocking_state,omitempty"`
	BlockingQuery    string     `gorm:"type:text" json:"blocking_query,omitempty"`
	State            string     `gorm:"size:64;index" json:"state"`
	Query            string     `gorm:"type:text" json:"query,omitempty"`
	QueryFingerprint string     `gorm:"size:64;index" json:"query_fingerprint,omitempty"`
	QueryDurationMS  int64      `json:"query_duration_ms"`
	ObservedAt       time.Time  `gorm:"not null;index" json:"observed_at"`
}
