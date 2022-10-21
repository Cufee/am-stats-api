package main

import (
	"github.com/byvko-dev/am-core/helpers/env"
	"github.com/byvko-dev/am-stats-api/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Setup a server
	app := fiber.New()

	app.Use(logger.New())

	v1 := app.Group("/v1")

	session := v1.Group("/session")
	session.Post("/player", handlers.GetPlayerSession)

	app.Listen(":" + env.MustGetString("PORT"))
}
