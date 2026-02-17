package repository

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/repository/mapper"
	"jello-api/pkg/couchdb"
	"time"
)

type IBookingRepository interface {
	GetActiveByDate(ctx context.Context, date time.Time) ([]domain.Booking, error)
	Create(ctx context.Context, table *domain.Booking) (*domain.Booking, error)
}

type couchBookingRepo struct {
	client *couchdb.Client
}

func NewBookingRepo(client *couchdb.Client) IBookingRepository {
	return &couchBookingRepo{client: client}
}

func (r *couchBookingRepo) GetActiveByDate(
	ctx context.Context,
	date time.Time,
) ([]domain.Booking, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	selector := map[string]interface{}{
		"type": "booking",
		"status": map[string]interface{}{
			"$in": []string{"checked_in", "reserved"},
		},
	}

	query := map[string]interface{}{
		"selector": selector,
	}

	rows, err := r.client.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find bookings: %w", err)
	}
	defer rows.Close()

	bookings, err := mapper.ScanAndMap(
		rows,
		mapper.BookingMapper{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map bookings: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}

func (r *couchBookingRepo) Create(ctx context.Context, d *domain.Booking) (*domain.Booking, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	doc := mapper.BookingMapper{}.ToModel(*d)
	doc.Type = "booking"

	if _, err := r.client.CreateDocWithID(ctx, doc.ID, doc); err != nil {
		return nil, err
	}

	return d, nil
}
