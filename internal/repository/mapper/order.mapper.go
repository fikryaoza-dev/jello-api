package mapper

import (
	"jello-api/internal/domain"
	"jello-api/internal/model"
)

type OrderMapper struct{}

func (o OrderMapper) ToDomain(m model.Order) domain.Order {
	return domain.Order{
		ID:           m.ID,
		BookingID:    m.BookingID,
		OrderNumber:  m.OrderNumber,
		CustomerName: m.CustomerName,
		TotalAmount:  m.TotalAmount,
		Items:        o.itemsToDomain(m.Items),
		Status:       m.Status,
	}
}

func (o OrderMapper) ToModel(d domain.Order) model.Order {
	return model.Order{
		ID:           d.ID,
		BookingID:    d.BookingID,
		OrderNumber:  d.OrderNumber,
		CustomerName: d.CustomerName,
		TotalAmount:  d.TotalAmount,
		CreatedAt:  d.CreatedAt,
		Status:       string(d.Status),
		Items:        o.itemsToModel(d.Items),
	}
}

func (OrderMapper) itemsToDomain(moItems []model.OrderItem) []domain.OrderItem {
	items := make([]domain.OrderItem, len(moItems))
	for i, item := range moItems {
		items[i] = domain.OrderItem{
			MenuID:   item.MenuID,
			MenuName: item.MenuName,
			Price:    item.Price,
			Quantity: item.Quantity,
			Subtotal: item.Subtotal,
			Notes:    item.Notes,
		}
	}
	return items
}

func (OrderMapper) itemsToModel(dItems []domain.OrderItem) []model.OrderItem {
	items := make([]model.OrderItem, len(dItems))
	for i, item := range dItems {
		items[i] = model.OrderItem{
			MenuID:   item.MenuID,
			MenuName: item.MenuName,
			Price:    item.Price,
			Quantity: item.Quantity,
			Subtotal: item.Subtotal,
			Notes:    item.Notes,
		}
	}
	return items
}
