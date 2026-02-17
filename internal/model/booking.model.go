package model

import (
	"time"
)

type Booking struct {
	ID              string    `json:"_id,omitempty"`
	Rev             string    `json:"_rev,omitempty"`
	Type            string    `json:"type"`
	Customer        string    `json:"customer"`
	TableID         string    `json:"table_id"`
	DurationMinutes int       `json:"duration_minutes"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
