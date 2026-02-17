package model

type Counter struct {
	ID    string `json:"_id"`
	Rev   string `json:"_rev,omitempty"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}
