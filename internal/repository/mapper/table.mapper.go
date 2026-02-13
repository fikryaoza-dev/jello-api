package mapper

import (
	"jello-api/internal/domain"
	"jello-api/internal/model"
)

func ToDomain(doc model.Table) domain.Table {
	return domain.Table{
		ID:     doc.ID,
		Name:   doc.Name,
		Status: domain.TableStatus(doc.Status),
	}
}

func ToModel(t domain.Table) model.Table {
	return model.Table{
		ID:     t.ID,
		Type:   "table",
		Name:   t.Name,
		Status: string(t.Status),
	}
}