package mapper

import (
	"jello-api/internal/domain"
	"jello-api/internal/model"
)

type TableMapper struct{}

func (TableMapper) ToDomain(m model.Table) domain.Table {
	return domain.Table{
		ID:       m.ID,
		Name:     m.Name,
		Capacity: m.Capacity,
		Area:     m.Area,
		Status:   domain.TableStatus(m.Status),
	}
}

func (TableMapper) ToModel(d domain.Table) model.Table {
	return model.Table{
		ID:       d.ID,
		Name:     d.Name,
		Capacity: d.Capacity,
		Area:     d.Area,
		Status:   string(d.Status),
	}
}
