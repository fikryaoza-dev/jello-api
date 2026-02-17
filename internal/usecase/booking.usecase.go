package usecase

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/dto"
	"jello-api/internal/repository"
	"math/rand"
	"time"
)

type BookingUsecase struct {
	Repo repository.IBookingRepository
}

func NewBookingUsecase(repo repository.IBookingRepository) *BookingUsecase {
	return &BookingUsecase{
		Repo: repo,
	}
}

func (u *BookingUsecase) CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (domain.Booking, error) {
	d := domain.Booking{
		ID:       GenerateBookingID(),
		Customer: req.CustomerName,
		TableID:  req.TableID,
		Status:   "checked_in",
	}
	_, err := u.Repo.Create(ctx, &d)
	if err != nil {
		return d, fmt.Errorf("failed to find menus: %w", err)
	}

	return d, nil
}

func GenerateBookingID() string {
	now := time.Now()
	random := rand.Intn(1000) // 0-999

	return fmt.Sprintf(
		"BK-%s-%03d",
		now.Format("20060102150405"),
		random,
	)
}
