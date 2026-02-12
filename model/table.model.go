package model

type Table struct {
	ID       string `json:"_id,omitempty"`
	Rev      string `json:"_rev,omitempty"`
	Type     string `json:"type"`
	Number   int    `json:"number"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}
