package main

import (
	"jello-api/config"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	database := config.ConnectDB()
	defer database.Client.Close()

	app := fiber.New(fiber.Config{
			AppName:      "Fiber CouchDB API v1.0.0",
			ServerHeader: "Fiber",
	})

		app.Use(recover.New())
	app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

		// Routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
					"status":      "ok",
					"message":     "Server is running",
					"environment": config.GetEnv(),
			})
	})

	port := config.GetPort()
    log.Printf("🚀 Server starting on port %s", port)
    log.Printf("📝 Environment: %s", config.GetEnv())
    log.Fatal(app.Listen(port))
}