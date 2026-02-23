package mapper

import (
	"jello-api/internal/domain"
	"jello-api/internal/model"
	"time"

	"github.com/google/uuid"
)

type OrderMapper struct{}

func (o OrderMapper) ToDomain(m model.Order) domain.Order {
	return domain.Order{
		ID:           m.ID,
		TableID:      m.TableID,
		Rev:          m.Rev,
		BookingID:    m.BookingID,
		OrderNumber:  m.OrderNumber,
		CustomerName: m.CustomerName,
		TotalAmount:  m.TotalAmount,
		Items:        o.itemsToDomain(m.Items),
		Status:       domain.OrderStatus(m.Status),
	}
}

func (o OrderMapper) ToModel(d domain.Order) model.Order {
	return model.Order{
		ID:           d.ID,
		BookingID:    d.BookingID,
		Type:         "order",
		Rev:          d.Rev,
		TableID:      d.TableID,
		OrderNumber:  d.OrderNumber,
		CustomerName: d.CustomerName,
		TotalAmount:  d.TotalAmount,
		CreatedAt:    d.CreatedAt,
		Status:       string(d.Status),
		Items:        o.itemsToModel(d.Items),
	}
}

func (OrderMapper) itemsToDomain(moItems []model.OrderItem) []domain.OrderItem {
	items := make([]domain.OrderItem, len(moItems))
	for i, item := range moItems {
		items[i] = domain.OrderItem{
			ID:       item.ID,
			MenuID:   item.MenuID,
			MenuName: item.MenuName,
			Price:    item.Price,
			Quantity: item.Quantity,
			Subtotal: item.Subtotal,
			Status:   domain.OrderItemStatus(item.Status),
			Notes:    item.Notes,
		}
	}
	return items
}

func (OrderMapper) itemsToModel(dItems []domain.OrderItem) []model.OrderItem {
	items := make([]model.OrderItem, len(dItems))
	for i, item := range dItems {
		items[i] = model.OrderItem{
			ID:        "order::" + uuid.New().String(),
			MenuID:    item.MenuID,
			MenuName:  item.MenuName,
			Price:     item.Price,
			Quantity:  item.Quantity,
			Subtotal:  item.Subtotal,
			Status:    string(item.Status),
			Notes:     item.Notes,
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		}
	}
	return items
}
