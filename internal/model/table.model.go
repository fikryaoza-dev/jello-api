package model

type Table struct {
	ID       string `json:"_id,omitempty"`
	Rev      string `json:"_rev,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Area     string `json:"area"`
	Number   int    `json:"number"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}
