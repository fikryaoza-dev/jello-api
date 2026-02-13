package response

import (
	"github.com/gofiber/fiber/v2"
)

type Meta struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type APIResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// SUCCESS RESPONSE
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(APIResponse{
		Meta: Meta{
			Success: true,
			Message: message,
			Error:   nil,
		},
		Data: data,
	})
}

// ERROR RESPONSE
func Fail(c *fiber.Ctx, code int, message string, err interface{}) error {
	return c.Status(code).JSON(APIResponse{
		Meta: Meta{
			Success: false,
			Message: message,
			Error:   err,
		},
		Data: nil,
	})
}