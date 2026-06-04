// internal/database/database.go
package database

import (
	"fmt"

	// GORM Imports
	"github.com/gabutlabs/godevopin/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB membuat dan mengembalikan koneksi ke database menggunakan GORM
func ConnectDB(cfg config.DBConfig) (*gorm.DB, error) {
	// Membuat DSN (Data Source Name) string
	stringHost := "host=%s user=%s dbname=%s port=%s sslmode=disable TimeZone=UTC"
	if cfg.Password != "" {
		stringHost += " password=" + cfg.Password
	}
	dsn := fmt.Sprintf(stringHost,
		cfg.Host, cfg.User, cfg.DBName, cfg.Port)
	// Membuka koneksi menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Successfully connected to the database using GORM!")
	return db, nil
}
