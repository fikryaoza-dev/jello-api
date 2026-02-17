package handler

import (
	"jello-api/internal/domain"
	"jello-api/internal/shared"
	"jello-api/internal/usecase"
	"jello-api/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TableHandler struct {
	Usecase *usecase.TableUsecase
}

func NewTableHandler(u *usecase.TableUsecase) *TableHandler {
	return &TableHandler{Usecase: u}
}

func (h *TableHandler) GetAllTables(c *fiber.Ctx) error {
	query := c.Queries()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := shared.NewPagination(page, limit)
	tables, total, err := h.Usecase.GetAllTables(c.Context(), query, pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	pagination.Calculate(total)
	return response.Success(*c, "Success Get List Tables", tables, &pagination)
}

func (h *TableHandler) CreateTable(c *fiber.Ctx) error {
	var table domain.Table

	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := h.Usecase.CreateTable(c.Context(), &table)
	if err != nil {
		return response.FiberErrorHandler(c, err)
	}

	return response.Success(*c, "Success Get List Tables", "", nil)
}

func (h *TableHandler) GetTableByID(c *fiber.Ctx) error {
	id := c.Params("id")
	tables, err := h.Usecase.GetTableByID(c.Context(), id)
	if err != nil {
		return response.FiberErrorHandler(c, err)
	}

	return response.Success(*c, "Success Get List Tables", tables, nil)
}

// func (h *TableHandler) GetTablesByStatus(c *fiber.Ctx) error {
// 	status := c.Query("status")
// 	if status == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "status is required",
// 		})
// 	}

// 	tables, err := h.Usecase.GetTablesByStatus(status)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(tables)
// }
