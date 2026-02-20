package dto

import "strings"

type CreateBookingRequest struct {
	CustomerName string `json:"customer_name" validate:"required,min=2"`
	TableID      string `json:"table_id" validate:"required"`
}

func (r *CreateBookingRequest) Normalize() {
	r.CustomerName = strings.ToUpper(r.CustomerName)
}