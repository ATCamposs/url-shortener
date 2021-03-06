package main

import (
	"database/sql"
	"log"

	"github.com/atcamposs/url-shortener/infrastructure"
	"github.com/atcamposs/url-shortener/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var DefaultConnection *sql.DB

// CreateServer creates a new Fiber instance
func CreateServer() *fiber.App {
	app := fiber.New()

	return app
}

func main() {
	//Connect to database
	infrastructure.StartPostgresConnection()
	infrastructure.DefaultConnection = *infrastructure.PostgresConnection

	app := CreateServer()

	app.Use(cors.New())

	router.SetupRoutes(app)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":3000"))
}
