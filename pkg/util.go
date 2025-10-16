package pkg

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Implementasi hashing password (misalnya menggunakan bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err // Ganti dengan hasil hashing sebenarnya
}

func CheckPasswordHash(password, hash string) bool {
	// Implementasi pengecekan password dengan hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateResponse(message string, data any) map[string]any {
	return map[string]any{"message": message, "data": data}
}

func GenerateErrorResponse(message string, errors any) map[string]any {
	return map[string]any{"message": message, "errors": errors}
}
