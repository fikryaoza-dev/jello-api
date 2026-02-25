package handler

import (
	"jello-api/internal/dto"
	"jello-api/internal/usecase"
	"log"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	Usecase *usecase.OrderUsecase
}

func NewOrderHandler(u *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{Usecase: u}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	booking, err := h.Usecase.CreateOrder(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create Order",
		"data":    booking,
	})
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order id is missing. ",
		})
	}

	var req dto.UpdateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	req.OrderID = orderID // inject from path param

	result, err := h.Usecase.UpdateOrder(c.Context(), req)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error update order. ",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create Order",
		"data":    result,
	})
}

func (h *OrderHandler) UpdateOrderItemNote(c *fiber.Ctx) error {
	var req dto.UpdateOrderItemNoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body. ",
		})
	}

	result, err := h.Usecase.UpdateOrderItemNote(c.Context(), req)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error update order item note. ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Order",
		"data":    result,
	})
}

func (h *OrderHandler) GetAllPendingOrderItems(c *fiber.Ctx) error {
	result, err := h.Usecase.GetAllPendingOrderItems(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error get pending order items. ",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Create Order",
		"data":    result,
	})
}

func (h *OrderHandler) GetAllReadyOrderItems(c *fiber.Ctx) error {
	result, err := h.Usecase.GetAllReadyOrderItems(c.Context())
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error get ready order items. ",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Get List Order",
		"data":    result,
	})
}

func (h *OrderHandler) GetAllServedOrderItems(c *fiber.Ctx) error {
	result, err := h.Usecase.GetAllServedOrderItems(c.Context())
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error get served order items. ",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Get List Order",
		"data":    result,
	})
}

func (h *OrderHandler) GetAllActiveOrders(c *fiber.Ctx) error {
	result, err := h.Usecase.GetActiveOrders(c.Context())
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "error get served order items. ",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Success Get List Order",
		"data":    result,
	})
}

func (h *OrderHandler) UpdateOrderItemStatus(c *fiber.Ctx) error {
	var req dto.UpdateOrderItemStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body. ",
		})
	}

	result, err := h.Usecase.UpdateOrderItemStatus(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body. ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success Update Order",
		"data":    result,
	})
}
