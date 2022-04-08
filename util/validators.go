package util

import (
	valid "github.com/asaskevich/govalidator"
	"github.com/atcamposs/url-shortener/models"
)

// IsEmpty checks if a string is empty
func IsEmpty(str string) (bool, string) {
	if valid.HasWhitespaceOnly(str) && str != "" {
		return true, "Must not be empty"
	}

	return false, ""
}

// ValidateRegister func validates the body of user for registration
func ValidateRegister(u *models.User) *models.UserErrors {
	e := &models.UserErrors{}
	e.Err, e.Username = IsEmpty(u.Username)

	if !valid.IsEmail(u.Email) {
		e.Err, e.Email = true, "Must be a valid email"
	}

	if len(u.Username) < 5 {
		e.Err, e.Username = true, "Username must have more than 5 characters"
	}

	//re := regexp.MustCompile(`\\d`) // regex check for at least one integer in string
	if !(len(u.Password) >= 8) {
		e.Err, e.Password = true, "Length of password should be atleast 8 and it must be a combination of uppercase letters, lowercase letters and numbers"
	}

	return e
}
