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
	// GetActiveByDate(ctx context.Context, date time.Time) ([]domain.Booking, error)
	Create(ctx context.Context, d *domain.Order) (*domain.Order, error)
}

type couchOrderRepo struct {
	client *couchdb.Client
}

func NewOrderRepo(client *couchdb.Client) IBookingRepository {
	return &couchBookingRepo{client: client}
}

func (r *couchBookingRepo) CreateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	now := time.Now().UTC()
	d.CreatedAt = now
	doc := mapper.OrderMapper{}.ToModel(*d)
	doc.Type = "order"

	rev, err := r.client.CreateDocWithID(ctx, doc.ID, doc)
    if err != nil {
        return nil, fmt.Errorf("failed to create order doc: %w", err)
    }
	d.Rev = rev
	return d, nil
}