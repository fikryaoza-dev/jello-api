package repository

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/model"
	"jello-api/internal/repository/mapper"
	"jello-api/pkg/couchdb"
	"time"
)

type IOrderRepository interface {
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	GetActiveOrders(ctx context.Context) ([]domain.Order, error)
	GetAllPendingOrderItems(ctx context.Context) ([]domain.OrderItem, error)
	GetAllReadyOrderItems(ctx context.Context) ([]domain.OrderItem, error)
	GetAllServedOrderItems(ctx context.Context) ([]domain.OrderItem, error)
	CreateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error)
	UpdateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error)
	UpdateOrderItemNote(ctx context.Context, orderID, menuID, note string) (*domain.Order, error)
	UpdateOrderItemStatus(ctx context.Context, orderID, menuID string, status domain.OrderItemStatus) (*domain.Order, error)
}

type couchOrderRepo struct {
	client *couchdb.Client
}

func NewOrderRepo(client *couchdb.Client) IOrderRepository {
	return &couchOrderRepo{client: client}
}

func (r *couchOrderRepo) UpdateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error) {
	// 1. Fetch the latest revision to avoid CouchDB conflict
	existing, err := r.GetOrderByID(ctx, d.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing order: %w", err)
	}
	// 2. Map domain → model, carry over immutable fields
	model := mapper.OrderMapper{}.ToModel(*d)
	model.Rev = existing.Rev                 // must match latest _rev
	model.Type = "order"                     // never mutate type
	model.OrderNumber = existing.OrderNumber // immutable
	model.TableID = existing.TableID         // immutable
	model.BookingID = existing.BookingID     // immutable
	model.CustomerName = existing.CustomerName
	model.CreatedAt = existing.CreatedAt

	// 3. Recalculate total from items
	var total float64
	for _, item := range model.Items {
		total += item.Subtotal
	}
	model.TotalAmount = total

	newRev, err := r.client.UpdateDocWithRev(ctx, model.ID, model.Rev, model)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// 5. Stamp the new revision and return
	model.Rev = newRev
	m := mapper.OrderMapper{}.ToDomain(*&model)
	return &m, nil
}

func (r *couchOrderRepo) CreateOrder(ctx context.Context, d *domain.Order) (*domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	orderId, err := r.client.GenerateOrderID(ctx)
	d.ID = orderId
	doc := mapper.OrderMapper{}.ToModel(*d)
	doc.Type = "order"

	rev, err := r.client.CreateDocWithID(ctx, doc.ID, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to create order doc: %w", err)
	}
	d.Rev = rev
	return d, nil
}

func (r *couchOrderRepo) GetActiveOrders(ctx context.Context) ([]domain.Order, error) {
	selector := map[string]interface{}{
		"type": "order",
		"status": map[string]interface{}{
			"$in": []string{"confirmed", "preparing", "served"},
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
		mapper.OrderMapper{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to map bookings: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bookings: %w", err)
	}

	return bookings, nil
}

func (r *couchOrderRepo) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	var model model.Order
	if err := r.client.GetDoc(ctx, id, &model); err != nil {
		return nil, fmt.Errorf("menu not found: %w", err)
	}
	menu := mapper.OrderMapper{}.ToDomain(model)
	return &menu, nil
}

func (r *couchOrderRepo) UpdateOrderItemNote(ctx context.Context, orderID, menuID, note string) (*domain.Order, error) {
	// 1. Fetch existing order
	existing, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	// 2. Find the item by menuID and update note
	found := false
	for i := range existing.Items {
		if existing.Items[i].MenuID == menuID {
			existing.Items[i].Notes = note
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// 3. Map to model
	m := mapper.OrderMapper{}.ToModel(*existing)
	// 4. Persist
	newRev, err := r.client.UpdateDocWithRev(ctx, m.ID, m.Rev, m)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item note: %w", err)
	}

	m.Rev = newRev
	model := mapper.OrderMapper{}.ToDomain(*&m)
	return &model, nil
}

func (r *couchOrderRepo) GetAllPendingOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	rows, err := r.client.Find(ctx, map[string]interface{}{
		"selector": map[string]interface{}{
			"type": "order",
			// "items": map[string]interface{}{
			// "$elemMatch": map[string]interface{}{
			// 	"status": domain.OrderItemStatusPending,
			// },
			// },
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var pendingItems []domain.OrderItem
	for rows.Next() {
		var m model.Order
		if err := rows.ScanDoc(&m); err != nil {
			return nil, fmt.Errorf("failed to scan order doc: %w", err)
		}
		for _, item := range m.Items {
			// if item.Status == string(domain.OrderItemStatusPending) {
			pendingItems = append(pendingItems, domain.OrderItem{
				OrderID:  m.ID,
				ID:       item.ID,
				MenuID:   item.MenuID,
				MenuName: item.MenuName,
				Price:    item.Price,
				Quantity: item.Quantity,
				Subtotal: item.Subtotal,
				Status:   domain.OrderItemStatus(item.Status),
				Notes:    item.Notes,
			})
			// }
		}
	}

	if len(pendingItems) == 0 {
		return nil, fmt.Errorf("no pending items found")
	}

	return pendingItems, nil
}

func (r *couchOrderRepo) UpdateOrderItemStatus(ctx context.Context, orderID, menuID string, status domain.OrderItemStatus) (*domain.Order, error) {
	// 1. Fetch existing order
	existing, err := r.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	// 2. Find item and update status
	found := false
	for i := range existing.Items {
		if existing.Items[i].MenuID == menuID {
			existing.Items[i].Status = status
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("menu item '%s' not found in order", menuID)
	}
	existing.Status = "served"
	// 3. Map to model and persist
	m := mapper.OrderMapper{}.ToModel(*existing)
	newRev, err := r.client.UpdateDocWithRev(ctx, m.ID, m.Rev, m)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item status: %w", err)
	}

	m.Rev = newRev
	model := mapper.OrderMapper{}.ToDomain(*&m)
	return &model, nil
}

func (r *couchOrderRepo) GetAllReadyOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	rows, err := r.client.Find(ctx, map[string]interface{}{
		"selector": map[string]interface{}{
			"type": "order",
			"items": map[string]interface{}{
				"$elemMatch": map[string]interface{}{
					"status": domain.OrderItemStatusReady,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var pendingItems []domain.OrderItem
	for rows.Next() {
		var m model.Order
		if err := rows.ScanDoc(&m); err != nil {
			return nil, fmt.Errorf("failed to scan order doc: %w", err)
		}
		for _, item := range m.Items {
			if item.Status == string(domain.OrderItemStatusReady) {
				pendingItems = append(pendingItems, domain.OrderItem{
					OrderID:  m.ID,
					ID:       item.ID,
					MenuID:   item.MenuID,
					TableID:  m.TableID,
					MenuName: item.MenuName,
					Price:    item.Price,
					Quantity: item.Quantity,
					Subtotal: item.Subtotal,
					Status:   domain.OrderItemStatus(item.Status),
					Notes:    item.Notes,
				})
			}
		}
	}

	// if len(pendingItems) == 0 {
	// 	return nil, fmt.Errorf("no ready items found")
	// }

	return pendingItems, nil
}

func (r *couchOrderRepo) GetAllServedOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	rows, err := r.client.Find(ctx, map[string]interface{}{
		"selector": map[string]interface{}{
			"type": "order",
			"items": map[string]interface{}{
				"$elemMatch": map[string]interface{}{
					"status": domain.OrderItemStatusServed,
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var pendingItems []domain.OrderItem
	for rows.Next() {
		var m model.Order
		if err := rows.ScanDoc(&m); err != nil {
			return nil, fmt.Errorf("failed to scan order doc: %w", err)
		}
		for _, item := range m.Items {
			if item.Status == string(domain.OrderItemStatusServed) {
				pendingItems = append(pendingItems, domain.OrderItem{
					OrderID:  m.ID,
					ID:       item.ID,
					MenuID:   item.MenuID,
					TableID:  m.TableID,
					MenuName: item.MenuName,
					Price:    item.Price,
					Quantity: item.Quantity,
					Subtotal: item.Subtotal,
					Status:   domain.OrderItemStatus(item.Status),
					Notes:    item.Notes,
				})
			}
		}
	}

	// if len(pendingItems) == 0 {
	// 	return nil, fmt.Errorf("no ready items found")
	// }
	return pendingItems, nil
}
