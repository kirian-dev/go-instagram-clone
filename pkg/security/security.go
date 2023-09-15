package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func DeletePassword(password *string) {
	if password != nil {
		*password = ""
	}
}

func ComparePasswords(hashedPassword, plaintextPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword)); err != nil {
		return err
	}
	return nil
}
