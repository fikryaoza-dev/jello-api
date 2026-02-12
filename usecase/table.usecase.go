package usecases

import (
	"context"
	"errors"
	"jello-api/config"
	"jello-api/model"
)

type TableUsecase struct {
	DB *config.Database
}

func NewTableUsecase(db *config.Database) *TableUsecase {
	return &TableUsecase{DB: db}
}

func (u *TableUsecase) CreateTable(table *model.Table) (string, string, error) {
	table.Type = "table"

	if table.Status == "" {
		table.Status = "available"
	}

	if table.Status != "available" && table.Status != "booked" && table.Status != "full" {
		return "", "", errors.New("invalid status")
	}

	ctx := context.Background()
	rev, err := u.DB.DB.Put(ctx, "", table)
	if err != nil {
		return "", "", err
	}

	return "", rev, nil
}

func (u *TableUsecase) GetAllTables() ([]model.Table, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"type": "table",
		},
	}

	ctx := context.Background()
	rows := u.DB.DB.Find(ctx, query)
	defer rows.Close()

	var tables []model.Table

	for rows.Next() {
		var t model.Table
		if err := rows.ScanDoc(&t); err == nil {
			tables = append(tables, t)
		}
	}

	return tables, nil
}

func (u *TableUsecase) GetTablesByStatus(status string) ([]model.Table, error) {
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"type":   "table",
			"status": status,
		},
	}

	ctx := context.Background()
	rows := u.DB.DB.Find(ctx, query)
	defer rows.Close()

	var tables []model.Table

	for rows.Next() {
		var t model.Table
		if err := rows.ScanDoc(&t); err == nil {
			tables = append(tables, t)
		}
	}

	return tables, nil
}
