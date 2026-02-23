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

	orderRepo := repository.NewOrderRepo(dbClient)
	tableRepo := repository.NewTableRepo(dbClient)
	tableUsecase := usecase.NewTableUsecase(tableRepo, bookingRepo, orderRepo)
	tableHandler := handler.NewTableHandler(tableUsecase)

	menuRepo := repository.NewMenuRepo(dbClient)
	menuUsecase := usecase.NewMenuUsecase(menuRepo, orderRepo)
	menuHandler := handler.NewMenuHandler(menuUsecase)

	orderUsecase := usecase.NewOrderUsecase(orderRepo, menuRepo)
	orderHandler := handler.NewOrderHandler(orderUsecase)

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

	order := app.Group("/api/v1/order")
	{
		order.Get("/items/pending", orderHandler.GetAllPendingOrderItems)
		order.Get("/items/status/ready", orderHandler.GetAllReadyOrderItems)
		order.Get("/items/status/served", orderHandler.GetAllServedOrderItems)
		order.Post("/", orderHandler.CreateOrder)
		order.Put("/:id", orderHandler.UpdateOrder)
		order.Put("/items/status", orderHandler.UpdateOrderItemStatus)
		order.Put("/items/note", orderHandler.UpdateOrderItemNote)
	}
}
