package domain

type TableStatus string

const (
	TableAvailable TableStatus = "available"
	TableOccupied  TableStatus = "occupied"
	TableReserved  TableStatus = "reserved"
)

type Table struct {
	ID       string `json:"id,omitempty"`
	Rev      string `json:"rev,omitempty"`
	Name     string `json:"name,omitempty"`
	Capacity int    `json:"capacity"`
	Area     string `json:"area"`
	Booking  *Booking `json:"booking"`
	Status   TableStatus `json:"status,omitempty"`
}
