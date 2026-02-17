package handler

import (
	"jello-api/internal/shared"
	"jello-api/internal/usecase"
	"jello-api/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	Usecase *usecase.MenuUsecase
}

func NewMenuHandler(u *usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{Usecase: u}
}

func (h *MenuHandler) GetAllMenu(c *fiber.Ctx) error {
	query := c.Queries()
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := shared.NewPagination(page, limit)
	tables, total, err := h.Usecase.GetAllMenu(c.Context(), query, pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	pagination.Calculate(total)
	return response.Success(*c, "Success Get List Menus", tables, &pagination)
}
