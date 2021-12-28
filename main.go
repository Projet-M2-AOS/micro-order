package main

import (
	"log"

	"github.com/b4cktr4ck5r3/micro-order/database"
	_ "github.com/b4cktr4ck5r3/micro-order/docs"
	"github.com/b4cktr4ck5r3/micro-order/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the root endpoint 😉",
		})
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("")
	router.OrderRoute(api.Group("/orders"))
}

// @title Order micro-service
// @version 1.0
// @description Order micro-service documentation.
func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	database.ConnectDB()

	setupRoutes(app)

	// port := config.Config("PORT")
	err := app.Listen(":9999")

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}
