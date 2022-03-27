package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yusufatalay/C1TakeHomeCase/database"
	"github.com/yusufatalay/C1TakeHomeCase/user"

	_ "github.com/yusufatalay/C1TakeHomeCase/database"
)

func registerRoutes(app *fiber.App) {
	app.Post("/api/v1/createuser", user.CreateUser)
}

func main() {
	app := fiber.New()

	// unhadler error here
	db, _ := database.DBConn.DB()
	defer db.Close()

	// Check if system is up and running, app.All will response to ALL HTTP methods with the sme way
	app.All("/health", func(c *fiber.Ctx) error {
		responsejson := struct {
			Message string `json:"message"`
			Type    bool   `json:"type"`
		}{
			Message: "System is healthy",
			Type:    true,
		}
		err := c.JSON(responsejson)
		return err
	})
	registerRoutes(app)
	err := app.Listen(":3000")

	if err != nil {
		log.Fatal("Cannot start the app: ", err)
	}

}
