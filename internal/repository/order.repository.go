package repository

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/repository/mapper"
	"jello-api/pkg/couchdb"
	"time"
)

type IOrderRepository interface {
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	CreateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error)
}

type couchOrderRepo struct {
	client *couchdb.Client
}

func NewOrderRepo(client *couchdb.Client) IOrderRepository {
	return &couchOrderRepo{client: client}
}

func (r *couchOrderRepo) CreateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	orderId, err := r.client.GenerateOrderID(ctx)
	d.ID = orderId
	doc := mapper.OrderMapper{}.ToModel(*d)
	doc.Type = "order"

	rev, err := r.client.CreateDocWithID(ctx, doc.ID, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create order doc: %w", err)
	}
	d.Rev = rev
	return d, nil
}

func (r *couchOrderRepo) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"booking_id": id,
		},
		"limit": 1,
	}

	rows, err := r.client.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find tables: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("order not found")
	}
	var order domain.Order

	if err := rows.ScanDoc(&order); err != nil {
		return nil, err
	}

	return &order, nil
}
