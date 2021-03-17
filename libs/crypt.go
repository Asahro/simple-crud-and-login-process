package libs

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}

	return true
}
