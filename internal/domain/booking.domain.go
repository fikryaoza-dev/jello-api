package domain

import "time"

type BookingStatus string

const (
	BookingReserved  BookingStatus = "reserved"
	BookingCheckedIn BookingStatus = "checked_in"
	BookingCompleted BookingStatus = "completed"
	BookingCancelled BookingStatus = "cancelled"
)

type Customer struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Notes string `json:"notes,omitempty"`
}

type Booking struct {
	ID              string        `json:"id"`
	Customer        string        `json:"customer"`
	TableID         string        `json:"tableId"`
	DurationMinutes int           `json:"durationMinutes"`
	Status          BookingStatus `json:"status"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}
