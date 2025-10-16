package service

import (
	"errors"
	"time"

	"github.com/gabutlabs/devopin/internal/config"
	"github.com/gabutlabs/devopin/internal/model"
	"github.com/gabutlabs/devopin/internal/repository"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	RegisterUser(username, email, password string) error
	GetUserProfile(userID uint) (*model.User, error)
	UpdateUserProfile(user *model.User) error
	DeleteUser(userID uint) error
	ListUsers() ([]model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	AuthenticateUser(email, password string) (map[string]any, error)
}
type userService struct {
	// Tambahkan dependensi yang diperlukan, misalnya database connection
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repository: repo}
}

// Implementasikan metode layanan sesuai kebutuhan
// Contoh: RegisterUser, AuthenticateUser, GetUserProfile, UpdateUserProfile, dll.
func (s *userService) RegisterUser(name, email, password string) error {
	// Logika untuk mendaftarkan user baru
	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return err
	}
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword, // Pastikan untuk mengenkripsi password sebelum menyimpannya
	}

	return s.repository.CreateUser(user)
}

func (s *userService) GetUserProfile(userID uint) (*model.User, error) {
	return s.repository.GetUserByID(userID)
}

func (s *userService) UpdateUserProfile(user *model.User) error {
	return s.repository.UpdateUser(user)
}

func (s *userService) DeleteUser(userID uint) error {
	return s.repository.DeleteUser(userID)
}

func (s *userService) ListUsers() ([]model.User, error) {
	return s.repository.ListUsers()
}

func (s *userService) GetUserByEmail(email string) (*model.User, error) {
	return s.repository.GetUserByEmail(email)
}

func (s *userService) AuthenticateUser(email string, password string) (map[string]any, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	user, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	matched := pkg.CheckPasswordHash(password, user.Password)
	if !matched {
		return nil, errors.New("invalid credentials")
	}
	claims := jwt.MapClaims{
		"id":    user.ID,
		"sub":   user.Email,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(cfg.AppSetting.JWTSecret))
	if err != nil {
		return nil, err
	}
	return map[string]any{"token": t, "user": user}, nil
}
