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
	"github.com/gofiber/fiber/v2"
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
		registerErrors.Err, registerErrors.Email = true, "Email is already registered"
		emailExists, _ := json.Marshal(registerErrors)
		return entity.User{}, errors.New(string(emailExists))
	}

	if userRepository.CheckUsernameExists(newUser.Username) {
		registerErrors.Err, registerErrors.Username = true, "Username is already registered"
		usernameExists, _ := json.Marshal(registerErrors)
		return entity.User{}, errors.New(string(usernameExists))
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
		registerErrors.Err, registerErrors.Username = true, "Something went wrong, please try again later. 😕"
		cannotRegister, _ := json.Marshal(registerErrors)
		return entity.User{}, errors.New(string(cannotRegister))
	}

	return user, nil

	// setting up the authorization cookies
	accessToken, refreshToken := util.GenerateTokens(user.UUID.String())
	accessCookie, refreshCookie := util.GetAuthCookies(accessToken, refreshToken)
	c.Cookie(accessCookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
