package model

import (
	"encoding/json"
	"time"
)

// Status untuk alarm yang sedang aktif
type ActiveAlarmStatus string

const (
	// Alarm sedang menyala dan akan mengirim notifikasi berulang
	StatusFiring ActiveAlarmStatus = "FIRING"
	// Alarm sedang menyala tapi sudah diakui, notifikasi berulang dihentikan
	StatusAcknowledged ActiveAlarmStatus = "ACKNOWLEDGED"
)

// ActiveAlarm menyimpan status LIVE dari alarm yang sedang menyala
type ActiveAlarm struct {
	// Primary Key gabungan untuk memastikan 1 alarm unik per target
	AlarmName string `gorm:"primarykey;type:varchar(100)" json:"alarm_name"`
	Target    string `gorm:"primarykey;type:varchar(255)" json:"target"`

	// Status saat ini: FIRING atau ACKNOWLEDGED
	Status ActiveAlarmStatus `gorm:"type:varchar(20);not null" json:"status"`

	// Kapan alarm ini pertama kali terdeteksi
	StartedAt time.Time `gorm:"not null" json:"started_at"`

	// Kapan terakhir notifikasi dikirim (untuk logika "Repeat Every")
	LastNotifiedAt time.Time `gorm:"not null" json:"last_notified_at"`

	// Siapa yang melakukan Acknowledge (jika sudah)
	AcknowledgedBy *string `json:"acknowledged_by,omitempty"`

	// Pesan error/status terakhir yang mudah dibaca
	Message string `gorm:"type:text" json:"message"`
}

type AlarmHistoryStatus string

const (
	// Kejadian saat alarm pertama kali menyala
	AlarmStatusFiring AlarmHistoryStatus = "FIRING"
	// Kejadian saat alarm pulih (kembali normal)
	AlarmStatusResolved AlarmHistoryStatus = "RESOLVED"
)

// AlarmHistory adalah tabel log yang mencatat setiap kejadian alarm
type AlarmHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	// Nama unik alarm (dari config, misal: "CPU_CRITICAL")
	AlarmName string `gorm:"type:varchar(100);not null;index" json:"alarm_name"`

	// Target dari alarm (misal: "worker:email_sender")
	Target string `gorm:"type:varchar(255);index" json:"target"`

	// Status kejadian: FIRING atau RESOLVED
	Status AlarmHistoryStatus `gorm:"type:varchar(20);not null" json:"status"`

	// Pesan yang mudah dibaca manusia
	Message string `gorm:"type:text" json:"message"`

	// (Opsional tapi direkomendasikan) Data mentah saat kejadian (value, threshold)
	Metadata json.RawMessage `json:"metadata"`
}