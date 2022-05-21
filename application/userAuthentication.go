package application

import (
	"encoding/json"
	"errors"

	"github.com/atcamposs/url-shortener/domain/user/entity"
	valueobject "github.com/atcamposs/url-shortener/domain/user/valueObject"
	"github.com/atcamposs/url-shortener/infrastructure/user/persistence"
	"github.com/atcamposs/url-shortener/provider/date"
	"github.com/atcamposs/url-shortener/provider/password"
	"github.com/atcamposs/url-shortener/util"
	"github.com/google/uuid"
)

var userRepository = new(persistence.PostgresRepository)
var passwordProvider = new(password.BcryptPassword)
var dateTimeProvider = new(date.TimeDate)

func Register(newUser *valueobject.NewUser) (entity.User, error) {
	// validate if the email, username and password are in correct format
	registerErrors := util.ValidateRegister(newUser)
	if registerErrors.Err {
		validationError, _ := json.Marshal(registerErrors)
		return entity.User{}, errors.New(string(validationError))
	}

	if userRepository.CheckEmailExists(newUser.Email) {
		return entity.User{}, errors.New("email is already registered")
	}

	if userRepository.CheckUsernameExists(newUser.Username) {
		return entity.User{}, errors.New("username is already registered")
	}

	now := dateTimeProvider.NowInRfc3339()
	user := entity.User{
		UUID:      uuid.New(),
		Email:     newUser.Email,
		Username:  newUser.Username,
		Password:  passwordProvider.Hash(newUser.Password),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if !userRepository.Create(user) {
		return entity.User{}, errors.New("cannot create user")
	}

	return user, nil
}

func Login(email string, password string) (entity.User, error) {
	user, loginError := userRepository.RetrieveUserByEmail(email)
	if loginError != nil {
		return entity.User{}, errors.New("user not found")
	}

	passwordMatch := passwordProvider.Match(user.Password, password)
	if !passwordMatch {
		return entity.User{}, errors.New("invalid password")
	}

	return user, nil
}
