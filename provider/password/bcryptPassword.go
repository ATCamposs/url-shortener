package password

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPassword struct {
}

func New() PasswordInterface {
	return &BcryptPassword{}
}

func (p *BcryptPassword) Hash(inputPassword string) string {
	password := []byte(inputPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword(
		password,
		rand.Intn(2),
	)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func (p *BcryptPassword) Match(inputPassword string, actualPassword string) bool {
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(inputPassword), []byte(actualPassword)); err != nil {
		return false
	}
	return true
}
