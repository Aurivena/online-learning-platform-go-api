package domain

import "golang.org/x/crypto/bcrypt"

const customCost = 14

func PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), customCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func PasswordVerify(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
