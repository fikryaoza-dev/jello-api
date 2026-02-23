package domain

type Order struct {
	ID           string      `json:"id"`             // CouchDB uses _id
	Rev          string      `json:"_rev,omitempty"` // CouchDB revision
	Type         string      `json:"type"`           // Always "order"
	TableID      string      `json:"table_id"`
	BookingID    string      `json:"booking_id"`
	OrderNumber  string      `json:"order_number"`
	CustomerName string      `json:"customer_name"`
	TotalAmount  float64     `json:"total_amount"`
	Status       OrderStatus `json:"status"`
	Items        []OrderItem `json:"items"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ID       string          `json:"id"`
	OrderID  string          `json:"order_id"`
	TableID  string          `json:"table_id"`
	MenuID   string          `json:"menu_id"`
	MenuName string          `json:"menu_name"`
	Price    float64         `json:"price"`
	Quantity int             `json:"quantity"`
	Subtotal float64         `json:"subtotal"`
	Status   OrderItemStatus `json:"status"`
	Notes    string          `json:"notes"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // Order placed, awaiting confirmation
	OrderStatusConfirmed OrderStatus = "confirmed" // Kitchen has accepted the order
	OrderStatusPreparing OrderStatus = "preparing" // Kitchen is actively preparing
	OrderStatusReady     OrderStatus = "ready"     // Ready to be served
	OrderStatusServed    OrderStatus = "served"    // Delivered to the table
	OrderStatusCompleted OrderStatus = "completed" // Paid and done
	OrderStatusCancelled OrderStatus = "cancelled" // Cancelled before completion
	OrderStatusOnHold    OrderStatus = "on_hold"   // Temporarily paused (e.g. item unavailable)
)

// IsValid checks if the status value is a known OrderStatus
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusConfirmed, OrderStatusPreparing,
		OrderStatusReady, OrderStatusServed, OrderStatusCompleted,
		OrderStatusCancelled, OrderStatusOnHold:
		return true
	}
	return false
}

type OrderItemStatus string

const (
	OrderItemStatusPending   OrderItemStatus = "pending"
	OrderItemStatusPreparing OrderItemStatus = "preparing"
	OrderItemStatusReady     OrderItemStatus = "ready"
	OrderItemStatusServed    OrderItemStatus = "served"
	OrderItemStatusCancelled OrderItemStatus = "cancelled"
)

func (s OrderStatus) IsTerminal() bool {
	return s == OrderStatusCompleted || s == OrderStatusCancelled
}

func (s OrderItemStatus) IsValid() bool {
	switch s {
	case OrderItemStatusPending, OrderItemStatusPreparing,
		OrderItemStatusReady, OrderItemStatusServed, OrderItemStatusCancelled:
		return true
	}
	return false
}
