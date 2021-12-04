package main

import (
	"log"

	"github.com/b4cktr4ck5r3/micro-order/config"
	"github.com/b4cktr4ck5r3/micro-order/database"
	"github.com/b4cktr4ck5r3/micro-order/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the root endpoint 😉",
		})
	})

	api := app.Group("")

	router.OrderRoute(api.Group("/order"))
}

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	database.ConnectDB()

	setupRoutes(app)

	port := config.Config("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}
