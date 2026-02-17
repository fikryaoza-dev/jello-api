// delivery/http/router.go
package route

import (
	"jello-api/internal/handler"
	"jello-api/internal/repository"
	"jello-api/internal/usecase"
	"jello-api/pkg/couchdb"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, dbClient *couchdb.Client) {
	bookingRepo := repository.NewBookingRepo(dbClient)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingUsecase)

	tableRepo := repository.NewTableRepo(dbClient)
	tableUsecase := usecase.NewTableUsecase(tableRepo, bookingRepo)
	tableHandler := handler.NewTableHandler(tableUsecase)

	menuRepo := repository.NewMenuRepo(dbClient)
	menuUsecase := usecase.NewMenuUsecase(menuRepo)
	menuHandler := handler.NewMenuHandler(menuUsecase)

	// Table routes
	tables := app.Group("/api/v1/tables")
	{
		tables.Get("/", tableHandler.GetAllTables)    // GET /api/v1/tables
		tables.Get("/:id", tableHandler.GetTableByID) // GET /api/v1/tables/:id
		tables.Post("/", tableHandler.CreateTable)    // POST /api/v1/tables
		// tables.Put("/:id", tableHandler.UpdateTable)    // PUT /api/v1/tables/:id
		// tables.Delete("/:id", tableHandler.DeleteTable) // DELETE /api/v1/tables/:id
	}

	menus := app.Group("/api/v1/menus")
	{
		menus.Get("/", menuHandler.GetAllMenu)
	}

	booking := app.Group("/api/v1/books")
	{
		booking.Post("/", bookingHandler.CreateBooking)
	}
}
