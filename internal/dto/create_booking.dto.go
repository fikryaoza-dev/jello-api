package dto

type CreateBookingRequest struct {
	CustomerName string `json:"customer_name" validate:"required,min=2"`
	TableID      string `json:"table_id" validate:"required"`
}
