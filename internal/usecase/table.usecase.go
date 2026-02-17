package usecase

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/repository"
	"jello-api/internal/shared"
	"log"
	"time"
)

type TableUsecase struct {
	Repo        repository.ITableRepository
	BookingRepo repository.IBookingRepository
}

func NewTableUsecase(repo repository.ITableRepository, bookRepo repository.IBookingRepository) *TableUsecase {
	return &TableUsecase{Repo: repo, BookingRepo: bookRepo}
}

func (u *TableUsecase) CreateTable(ctx context.Context, table *domain.Table) error {
	err := u.Repo.Create(ctx, table)
	if err != nil {
		return err
	}

	return nil
}

func (u *TableUsecase) GetAllTables(ctx context.Context, queries map[string]string, pagination shared.Pagination) ([]domain.Table, int, error) {
	rows, total, err := u.Repo.GetAll(ctx, queries, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find tables: %w", err)
	}
	// 3️⃣ Ambil booking aktif di tanggal itu
	now := time.Now().UTC()
	today := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid date format")
	}
	bookings, err := u.BookingRepo.GetActiveByDate(ctx, today)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings: %w", err)
	}
	bookingMap := make(map[string]bool)
	for _, b := range bookings {
		bookingMap[b.TableID] = true
	}
	for i := range rows {
		log.Println(rows[i].ID, bookingMap[rows[i].ID])
		if bookingMap[rows[i].ID] {
			rows[i].Status = "available"
		} else {
			rows[i].Status = "booked"
		}
	}
	return rows, total, nil
}

func (u *TableUsecase) GetTableByID(ctx context.Context, id string) (*domain.Table, error) {
	// Get table from repository
	table, err := u.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find tables: %w", err)
	}

	return table, nil
}
