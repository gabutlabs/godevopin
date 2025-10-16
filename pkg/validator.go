package pkg

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// GetValidator tetap menjadi dasar/mesin dari singleton
func GetValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// ValidateAndFormat adalah FUNGSI BARU yang akan Anda panggil dari handler.
// Ia mengembalikan nil jika valid, atau map[string]string jika ada error.
func ValidateAndFormat(s interface{}) map[string]string {
	// Panggil validator
	err := GetValidator().Struct(s)

	// Jika tidak ada error, kembalikan nil (sukses)
	if err == nil {
		return nil
	}

	// Jika ada error, ubah menjadi ValidationErrors
	validationErrors := err.(validator.ValidationErrors)

	// Buat map untuk menampung error yang sudah diformat
	errors := make(map[string]string)

	// Iterasi error dan buat pesan yang lebih ramah
	for _, e := range validationErrors {
		// Anda bisa membuat pesan yang lebih spesifik di sini
		errors[e.Field()] = generateValidationMessage(e)
	}

	return errors
}

// generateValidationMessage adalah helper internal untuk membuat pesan error
func generateValidationMessage(e validator.FieldError) string {
	// Pesan error bisa dibuat lebih dinamis dan ramah
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Field %s is required.", e.Field())
	case "email":
		return fmt.Sprintf("Field %s must be a valid email address.", e.Field())
	case "min":
		return fmt.Sprintf("Field %s must be at least %s characters long.", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("Field %s must be at most %s characters long.", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("Field %s must be greater than or equal to %s.", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("Field %s must be less than or equal to %s.", e.Field(), e.Param())
	case "unique":
		return fmt.Sprintf("Field %s already exists", e.Field())
	default:
		return fmt.Sprintf("Field %s is not valid.", e.Field())
	}
}
