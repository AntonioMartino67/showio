package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword genera un hash sicuro della password in chiaro
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword confronta una password in chiaro con un hash salvato
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}