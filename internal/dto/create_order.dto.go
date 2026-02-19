package dto

type CreateOrderRequest struct {
	BookingID string             `json:"booking_id" validate:"required"`
	Items     []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateOrderRequest struct {
	Items  []OrderItemRequest `json:"items,omitempty" validate:"omitempty,dive"`
	Status string             `json:"status,omitempty"`
}

type OrderItemRequest struct {
	MenuID   string `json:"menu_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
	Notes    string `json:"notes,omitempty"`
}
