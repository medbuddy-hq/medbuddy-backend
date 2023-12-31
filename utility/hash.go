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
	space := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRST1234567890"
	var salt string

	random := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

	for i := 0; i < saltLen; i++ {
		idx := random.Intn(len(space))
		salt += string(space[idx])
	}

	return salt
}
