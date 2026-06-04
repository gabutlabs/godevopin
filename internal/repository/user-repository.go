package repository

import (
	"github.com/gabutlabs/godevopin/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	ListUsers() ([]model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	GetUserByEmail(email string) (*model.User, error)
}

type userRepository struct {
	// Tambahkan field yang diperlukan, misalnya koneksi database
	db *gorm.DB
}

// NewUserRepository membuat instance baru dari UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Implementasikan metode CRUD sesuai kebutuhan
// Contoh: CreateUser, GetUserByID, UpdateUser, DeleteUser, dll.
// Contoh metode untuk mendapatkan user berdasarkan ID
func (r *userRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
func (r *userRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
