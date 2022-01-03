package main

import (
	"log"

	"github.com/b4cktr4ck5r3/micro-order/config"
	"github.com/b4cktr4ck5r3/micro-order/database"
	_ "github.com/b4cktr4ck5r3/micro-order/docs"
	"github.com/b4cktr4ck5r3/micro-order/handler"
	"github.com/b4cktr4ck5r3/micro-order/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// Add routes to fiber instance
func setupRoutes(app *fiber.App) {
	//swagger doc
	app.Get("/swagger/*", fiberSwagger.WrapHandler)
	//swagger json format
	app.Get("/docs-json", handler.GetSwaggerJson)

	//orders endpoints
	api := app.Group("")
	router.OrderRoute(api.Group("/orders"))
}

// @title Order micro-service
// @version 1.0
// @description Order micro-service documentation.
func main() {
	//create new fiber instance
	app := fiber.New()

	//allow cors
	app.Use(cors.New())
	//allow logs
	app.Use(logger.New())

	//connect to mongo database
	database.ConnectDB()

	//add routes
	setupRoutes(app)

	//get port from env variable
	port := config.Config("PORT")
	//run web server instance
	err := app.Listen(":" + port)

	//check for error during running
	if err != nil {
		log.Fatal("Error app failed to start : " + err.Error())
		panic(err)
	}
}
