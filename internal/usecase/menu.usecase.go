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
}

func NewMenuUsecase(repo repository.IMenuRepository) *MenuUsecase {
	return &MenuUsecase{
		Repo: repo,
	}
}

func (u *MenuUsecase) GetAllMenu(ctx context.Context, queries map[string]string, pagination shared.Pagination) ([]domain.Menu, int, error) {
	rows, total, err := u.Repo.GetAll(ctx, queries, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find menus: %w", err)
	}

	return rows, total, nil
}
