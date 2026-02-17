package domain

import (
	"time"
)

type MenuStatus string

const (
	MenuAvailable MenuStatus = "available"
	MenuOutStock  MenuStatus = "out_of_stock"
	MenuInactive  MenuStatus = "inactive"
)

type Menu struct {
	ID          string     `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Price       float64    `json:"price"`
	Status      MenuStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
