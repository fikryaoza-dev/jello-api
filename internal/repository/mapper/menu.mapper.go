package mapper

import (
	"jello-api/internal/domain"
	"jello-api/internal/model"
)

type MenuMapper struct{}

func (MenuMapper) ToDomain(m model.Menu) domain.Menu {
	return domain.Menu{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Price:       m.Price,
		Status:      domain.MenuStatus(m.Status),
	}
}

func (MenuMapper) ToModel(d domain.Menu) model.Menu {
	return model.Menu{
		ID:          d.ID,
		Code:        d.Code,
		Name:        d.Name,
		Description: d.Description,
		Category:    d.Category,
		Price:       d.Price,
		Status:      string(d.Status),
	}
}
