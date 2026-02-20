package usecase

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/repository"
	"jello-api/internal/shared"
)

type MenuUsecase struct {
	Repo repository.IMenuRepository
	OrderRepo repository.IOrderRepository
}

func NewMenuUsecase(repo repository.IMenuRepository, orderRepo repository.IOrderRepository) *MenuUsecase {
	return &MenuUsecase{
		Repo: repo,
		OrderRepo: orderRepo,
	}
}

func (u *MenuUsecase) GetAllMenu(ctx context.Context, queries map[string]string, pagination shared.Pagination) ([]domain.Menu, int, error) {
	rows, total, err := u.Repo.GetAll(ctx, queries, pagination)
	if status, ok := queries["booking_id"]; ok && status != "" {
		order, err := u.OrderRepo.GetOrderByID(ctx, queries["booking_id"])
		if err != nil {
			return nil, 0, fmt.Errorf("failed to find order: %w", err)
		}
		for i := range rows {
			for _, item := range order.Items {
				if rows[i].ID == item.MenuID {
					rows[i].Quantity = int64(item.Quantity)
					break
				}
			}
		}
	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find menus: %w", err)
	}
	return rows, total, nil
}
