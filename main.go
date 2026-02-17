package main

import (
	"context"
	"jello-api/config"
	"jello-api/internal/route"
	"jello-api/pkg/couchdb"
	"jello-api/response"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	} else {
		log.Println(".env loaded successfully")
	}
	ctx := context.Background()
	dbClient, err := couchdb.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to CouchDB: %v", err)
	}
	defer dbClient.Client.Close()

	app := fiber.New(fiber.Config{
		AppName:      "Jello POS API v1.0.0",
		ServerHeader: "Fiber",
		ErrorHandler: response.FiberErrorHandler,
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
	route.SetupRoutes(app, dbClient)
	// Routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "ok",
			"message":     "Server is running",
			"environment": config.GetEnv("APP_ENV", ""),
		})
	})
	// api.Get("/tables/status", tableHandler.GetTablesByStatus)

	// PRINT ROUTES
	log.Println("========== ROUTES ==========")
	for _, r := range app.GetRoutes() {
		log.Printf("➡️  %-6s | %s", r.Method, r.Path)
	}
	log.Println("============================")
	port := config.GetPort()
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📝 Environment: %s", config.GetEnv("APP_ENV", ""))
	log.Fatal(app.Listen(port))
}
