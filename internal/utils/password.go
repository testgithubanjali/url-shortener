package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword converts a plain password into a bcrypt hash.
func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain password.
func CheckPassword(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
