package model

import (
	"time"
)

type Menu struct {
	ID          string    `json:"_id,omitempty"`
	Rev         string    `json:"_rev,omitempty"`
	Type        string    `json:"type"` // for filtering in CouchDB
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
