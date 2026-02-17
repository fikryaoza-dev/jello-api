package shared

import "log"

type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalRows  int  `json:"total_rows"`
	TotalPages int  `json:"total_pages"`
	NextPage   bool `json:"next_page"`
	PrevPage   bool `json:"prev_page"`
}

func NewPagination(page, limit int) Pagination {
	log.Println(page, limit)
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return Pagination{
		Page:  page,
		Limit: limit,
	}
}

func (p Pagination) Skip() int {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) Calculate(total int) {
	p.TotalRows = total

	if p.Limit <= 0 {
		p.Limit = 10
	}

	p.TotalPages = (total + p.Limit - 1) / p.Limit

	// Next Page
	if p.Page < p.TotalPages {
		// next := p.Page + 1
		p.NextPage = true
	}

	// Previous Page
	if p.Page > 1 {
		// prev := p.Page - 1
		p.PrevPage = true
	}
}
