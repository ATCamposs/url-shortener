package router

import (
	"os"
	"time"

	"github.com/atcamposs/url-shortener/application"
	db "github.com/atcamposs/url-shortener/database"
	"github.com/atcamposs/url-shortener/domain/user/entity"
	valueobject "github.com/atcamposs/url-shortener/domain/user/valueObject"
	"github.com/atcamposs/url-shortener/infrastructure/user/persistence"
	"github.com/atcamposs/url-shortener/provider/date"
	"github.com/atcamposs/url-shortener/provider/password"
	"github.com/atcamposs/url-shortener/util"
	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
)

var jwtKey = []uint8(os.Getenv("PRIV_KEY"))

var passwordProvider = new(password.BcryptPassword)
var dateProvider = new(date.TimeDate)
var userRepository = new(persistence.PostgresRepository)

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	USER.Post("/signup", CreateUser)              // Sign Up a user
	USER.Post("/signin", LoginUser)               // Sign In a user
	USER.Get("/get-access-token", GetAccessToken) // returns a new access_token

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(util.SecureAuth()) // middleware to secure all routes for this group
	privUser.Get("/user", GetUserData)
}

// CreateUser route registers a User into the database
func CreateUser(c *fiber.Ctx) error {
	newUser := new(valueobject.NewUser)

	if err := c.BodyParser(newUser); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	user, registrationError := application.Register(newUser)
	if registrationError != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"value": registrationError,
		})
	}

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

// LoginUser route logins a user in the app
func LoginUser(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input := new(LoginInput)

	if err := c.BodyParser(input); err != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"input": "Please review your input",
		})
	}

	user, loginError := application.Login(input.Email, input.Password)
	if loginError != nil {
		return c.JSON(fiber.Map{
			"error": true,
			"user":  "invalid email or password",
		})
	}

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

// GetUserData returns the details of the user signed in
func GetUserData(c *fiber.Ctx) error {
	id := c.Locals("id")

	u := new(entity.User)
	if res := db.DB.Where("uuid = ?", id).First(&u); res.RowsAffected <= 0 {
		return c.JSON(fiber.Map{"error": true, "general": "Cannot find the User"})
	}

	return c.JSON(u)
}

// GetAccessToken generates and sends a new access token iff there is a valid refresh token
func GetAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	refreshClaims := new(entity.Claims)
	token, _ := jwt.ParseWithClaims(refreshToken, refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if res := db.DB.Where(
		"expires_at = ? AND issued_at = ? AND issuer = ?",
		refreshClaims.ExpiresAt, refreshClaims.IssuedAt, refreshClaims.Issuer,
	).First(&entity.Claims{}); res.RowsAffected <= 0 {
		// no such refresh token exist in the database
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			// refresh token is expired
			c.ClearCookie("access_token", "refresh_token")
			return c.SendStatus(fiber.StatusForbidden)
		}
	} else {
		// malformed refresh token
		c.ClearCookie("access_token", "refresh_token")
		return c.SendStatus(fiber.StatusForbidden)
	}

	_, accessToken := util.GenerateAccessClaims(refreshClaims.Issuer)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(fiber.Map{"access_token": accessToken})
}
