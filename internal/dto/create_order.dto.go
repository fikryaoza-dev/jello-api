package dto

import "jello-api/internal/domain"

type CreateOrderRequest struct {
	TableID      string `json:"table_id" validate:"required"`
	CustomerName string `json:"customer_name" validate:"required"`
}

type UpdateOrderRequest struct {
	OrderID string             `json:"order_id" validate:"required"`
	Items   []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateOrderItemNoteRequest struct {
	OrderID string `json:"order_id" validate:"required"`
	MenuID  string `json:"menu_id"  validate:"required"`
	Note    string `json:"note"`
}

type OrderItemRequest struct {
	MenuID   string `json:"menu_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
	Notes    string `json:"notes,omitempty"`
}

type UpdateOrderItemStatusRequest struct {
	OrderID string                 `json:"order_id" validate:"required"`
	MenuID  string                 `json:"menu_id"  validate:"required"`
	Status  domain.OrderItemStatus `json:"status"   validate:"required,oneof=pending preparing ready served cancelled"`
}
