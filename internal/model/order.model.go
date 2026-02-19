package model

type Order struct {
	ID           string      `json:"_id,omitempty"`
	Rev          string      `json:"_rev,omitempty"`
	Type         string      `json:"type"` // Hardcoded to "order"
	BookingID    string      `json:"booking_id"`
	OrderNumber  string      `json:"order_number"`
	CustomerName string      `json:"customer_name"`
	TotalAmount  float64     `json:"total_amount"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items"`
	CreatedAt    string      `json:"created_at"` // CouchDB prefers ISO8601 strings
}

type OrderItem struct {
	MenuID   string  `json:"menu_id"`
	MenuName string  `json:"menu_name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
	Notes    string  `json:"notes"`
}
