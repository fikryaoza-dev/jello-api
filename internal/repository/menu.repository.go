package repository

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/repository/mapper"
	"jello-api/internal/shared"
	"jello-api/pkg/couchdb"
	"time"
)

type IMenuRepository interface {
	GetAll(ctx context.Context, filter map[string]string, pagination shared.Pagination) ([]domain.Menu, int, error)
}

type couchMenuRepo struct {
	client *couchdb.Client
}

func NewMenuRepo(client *couchdb.Client) IMenuRepository {
	return &couchMenuRepo{client: client}
}

// GetAll retrieves all tables
func (r *couchMenuRepo) GetAll(ctx context.Context, filter map[string]string, pagination shared.Pagination) ([]domain.Menu, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// base selector
	selector := map[string]interface{}{
		"type": "menu",
	}
	if status, ok := filter["status"]; ok && status != "" {
		selector["status"] = status
	}

	if keyword, ok := filter["keyword"]; ok && keyword != "" {
		selector["name"] = map[string]interface{}{
			"$regex": "(?i).*" + keyword + ".*",
		}
	}

	query := map[string]interface{}{
		"selector": selector,
		"limit":    pagination.Limit,
		"skip":     pagination.Skip(), // helper dari struct
	}

	rows, err := r.client.Find(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find tables: %w", err)
	}
	defer rows.Close()

	tables, err := mapper.ScanAndMap(
		rows,
		mapper.MenuMapper{},
	)

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating tables: %w", err)
	}

	countRows, err := r.client.Find(ctx, map[string]interface{}{
		"selector": selector,
	})
	if err != nil {
		return nil, 0, err
	}
	defer countRows.Close()
	total := 0
	for countRows.Next() {
		total++
	}
	return tables, total, nil
}
