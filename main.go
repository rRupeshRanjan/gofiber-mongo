package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gofiber-mongo/services"
	"log"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	registerRoutes(app)

	err := app.Listen(":8080")
	if err != nil {
		log.Panic("Error starting server with error: ", err)
	}
}

func registerRoutes(app *fiber.App) {
	app.Get("/book/:id", services.GetBookByIdHandler)
	app.Post("/book", services.CreateBookHandler)
	app.Put("/book/:id", services.UpdateBookHandler)
	app.Delete("/book/:id", services.DeleteBookHandler)
}
