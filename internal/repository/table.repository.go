package repository

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/model"
	"jello-api/internal/repository/mapper"
	"jello-api/internal/shared"
	"jello-api/pkg/couchdb"
	"time"
)

type ITableRepository interface {
	Create(ctx context.Context, table *domain.Table) error
	GetAll(ctx context.Context, filter map[string]string, pagination shared.Pagination) ([]domain.Table, int, error)
	GetByID(ctx context.Context, id string) (*domain.Table, error)
}

type couchTableRepo struct {
	client *couchdb.Client
}

func NewTableRepo(client *couchdb.Client) ITableRepository {
	return &couchTableRepo{client: client}
}

// GetAll retrieves all tables
func (r *couchTableRepo) GetAll(ctx context.Context, filter map[string]string, pagination shared.Pagination) ([]domain.Table, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	// base selector
	selector := map[string]interface{}{
		"type": "table",
	}
	if status, ok := filter["status"]; ok && status != "" {
		selector["status"] = status
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
		mapper.TableMapper{},
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

func (r *couchTableRepo) GetByID(ctx context.Context, id string) (*domain.Table, error) {
	var doc model.Table
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"_id": id,
		},
		"limit": 1,
	}
	// Get full document including revision
	err := r.client.FindOne(ctx, query, &doc)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	tableMapper := mapper.TableMapper{}
	// Convert to domain with revision
	table := tableMapper.ToDomain(doc)
	// table.Rev now contains current revision

	return &table, nil
}

// Create creates a new table document
func (r *couchTableRepo) Create(ctx context.Context, table *domain.Table) error {
	// // Convert domain to model

	// if table.Status == "" {
	// 	table.Status = "available"
	// }

	// if table.Status != "available" && table.Status != "booked" && table.Status != "full" {
	// 	return errors.New("invalid status")
	// }
	// doc := mapper.ToModel(*table)

	// // Create document with ID
	// rev, err := r.client.CreateDocWithID(ctx, doc.ID, doc)
	// if err != nil {
	// 	return fmt.Errorf("failed to create table: %w", err)
	// }

	// // Update domain object with revision
	// table.Rev = rev

	return nil
}

// Update updates an existing table document
// func (r *couchTableRepo) Update(ctx context.Context, table *domain.Table) error {
// 	// Validate revision exists
// 	if table.Rev == "" {
// 		return errors.New("revision required for update")
// 	}

// 	// Validate status
// 	if table.Status != "available" && table.Status != "booked" && table.Status != "full" {
// 		return errors.New("invalid status")
// 	}

// 	// Convert domain to model
// 	doc := mapper.ToModel(*table)

// 	// Update document
// 	rev, err := r.client.UpdateDoc(ctx, doc.ID, doc.Rev, doc)
// 	if err != nil {
// 		return fmt.Errorf("failed to update table: %w", err)
// 	}

// 	// Update domain object with new revision
// 	table.Rev = rev

// 	return nil
// }
