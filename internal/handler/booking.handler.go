package handler

import (
	"jello-api/internal/dto"
	"jello-api/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	Usecase *usecase.BookingUsecase
}

func NewBookingHandler(u *usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{Usecase: u}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req dto.CreateBookingRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	booking, err := h.Usecase.CreateBooking(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create Booking",
		"data":    booking,
	})
}
