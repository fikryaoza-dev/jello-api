package domain

type TableStatus string

const (
	TableAvailable TableStatus = "available"
	TableOccupied  TableStatus = "occupied"
	TableReserved  TableStatus = "reserved"
)

type Table struct {
	ID     string
	Name   string
	Status TableStatus
}