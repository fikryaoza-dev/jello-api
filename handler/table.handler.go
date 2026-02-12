package handler

import (
	"jello-api/model"
	usecases "jello-api/usecase"

	"github.com/gofiber/fiber/v2"
)

type TableHandler struct {
	Usecase *usecases.TableUsecase
}

func NewTableHandler(u *usecases.TableUsecase) *TableHandler {
	return &TableHandler{Usecase: u}
}

func (h *TableHandler) CreateTable(c *fiber.Ctx) error {
	var table model.Table

	if err := c.BodyParser(&table); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, rev, err := h.Usecase.CreateTable(&table)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":  id,
		"rev": rev,
	})
}

func (h *TableHandler) GetAllTables(c *fiber.Ctx) error {
	tables, err := h.Usecase.GetAllTables()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tables)
}

func (h *TableHandler) GetTablesByStatus(c *fiber.Ctx) error {
	status := c.Query("status")
	if status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status is required",
		})
	}

	tables, err := h.Usecase.GetTablesByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(tables)
}
