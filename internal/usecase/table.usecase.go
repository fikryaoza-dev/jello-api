package usecase

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/repository"
	"jello-api/internal/shared"
)

type TableUsecase struct {
	Repo        repository.ITableRepository
	BookingRepo repository.IBookingRepository
	OrderRepo   repository.IOrderRepository
}

func NewTableUsecase(repo repository.ITableRepository, bookRepo repository.IBookingRepository, orderRepo repository.IOrderRepository) *TableUsecase {
	return &TableUsecase{Repo: repo, BookingRepo: bookRepo, OrderRepo: orderRepo}
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
	if err != nil {
		return nil, 0, fmt.Errorf("invalid date format")
	}
	orders, err := u.OrderRepo.GetActiveOrders(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get bookings: %w", err)
	}
	orderMap := make(map[string]domain.Order, len(orders))
	for _, b := range orders {
		orderMap[b.TableID] = domain.Order{
			ID:           b.ID,
			Type:         b.Type,
			CustomerName: b.CustomerName,
			OrderNumber:  b.OrderNumber,
			Status:       b.Status,
			TableID:      b.TableID,
			Items:        b.Items,
		}
	}
	for i := range rows {
		if detail, exists := orderMap[rows[i].ID]; exists {
			rows[i].Status = "occupied"
			rows[i].Order = &detail // Assign the pointer to the struct
		} else {
			rows[i].Status = "available"
			rows[i].Order = nil // Ensure it's null if not booked
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
