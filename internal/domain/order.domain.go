package domain

import "time"

type Order struct {
	ID           string      `json:"_id"`            // CouchDB uses _id
	Rev          string      `json:"_rev,omitempty"` // CouchDB revision
	Type         string      `json:"type"`           // Always "order"
	BookingID    string      `json:"booking_id"`
	OrderNumber  string      `json:"order_number"`
	CustomerName string      `json:"customer_name"`
	TotalAmount  float64     `json:"total_amount"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items"`
	CreatedAt    time.Time   `json:"created_at"`
}

type OrderItem struct {
	MenuID   string  `json:"menu_id"`
	MenuName string  `json:"menu_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
	Notes    string  `json:"notes"`
}
