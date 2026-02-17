package response

import (
	"jello-api/internal/shared"

	"github.com/gofiber/fiber/v2"
)

type Meta struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type APIResponse[T any] struct {
	Meta       Meta               `json:"meta"`
	Data       T                  `json:"data"`
	Pagination *shared.Pagination `json:"pagination,omitempty"`
}

type APIPaginationResponse[T any] struct {
	Meta       Meta               `json:"meta"`
	Data       T                  `json:"data"`
	Pagination *shared.Pagination `json:"pagination,omitempty"`
}

func Success[T any](c fiber.Ctx, message string, data T, pagination *shared.Pagination) error {
	return c.JSON(APIPaginationResponse[T]{
		Meta: Meta{
			Success: true,
			Message: message,
			Error:   nil,
		},
		Data:       data,
		Pagination: pagination,
	})
}
