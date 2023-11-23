package utility

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

var saltLen = 8

func HashPassword(password string) (hashed string, salt string, err error) {
	salt = randomSalt()
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return string(hash), salt, nil
}

func PasswordIsValid(password, salt, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	return err == nil
}

func randomSalt() string {
	var salt string
	random := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < saltLen; i++ {
		char := random.Int31n(26) + 61
		salt += string(char)
	}

	return salt
}
