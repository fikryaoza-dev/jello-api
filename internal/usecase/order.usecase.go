package usecase

import (
	"context"
	"fmt"
	"jello-api/internal/domain"
	"jello-api/internal/dto"
	"jello-api/internal/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type OrderUsecase struct {
	Repo     repository.IOrderRepository
	MenuRepo repository.IMenuRepository
}

func NewOrderUsecase(repo repository.IOrderRepository, menuRepo repository.IMenuRepository) *OrderUsecase {
	return &OrderUsecase{
		Repo:     repo,
		MenuRepo: menuRepo,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*domain.Order, error) {
	// 2. Build the order domain object
	now := time.Now().UTC()
	order := &domain.Order{
		TableID:      req.TableID,      // e.g. "order::" + uuid
		CustomerName: req.CustomerName, // e.g. "order::" + uuid
		Type:         "order",
		OrderNumber:  generateOrderNumber(now), // e.g. "ORD-20240315-0001"
		Status:       "confirmed",
		CreatedAt:    NowWIB(),
	}

	// 3. Persist
	result, err := u.Repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return result, nil
}

func (u *OrderUsecase) UpdateOrder(ctx context.Context, req dto.UpdateOrderRequest) (*domain.Order, error) {
	// 1. Fetch existing order to ensure it exists and is updatable
	existing, err := u.Repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// 2. Guard against updating a terminal order
	if existing.Status.IsTerminal() {
		return nil, fmt.Errorf("cannot update order with status %q", existing.Status)
	}

	// 3. Resolve and validate each menu item from request
	items, total, err := u.resolveOrderItems(ctx, *existing, req.Items)
	if err != nil {
		return nil, err
	}

	// 4. Build updated domain object — only mutate what's allowed
	existing.Items = items
	existing.TotalAmount = total

	// 5. Persist
	result, err := u.Repo.UpdateOrder(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return result, nil
}

func (u *OrderUsecase) resolveOrderItems(ctx context.Context, orderDet domain.Order, reqItems []dto.OrderItemRequest) ([]domain.OrderItem, float64, error) {
	items := orderDet.Items
	total := 0.0

	existing := make(map[string]*domain.OrderItem)
	for i := range items {
		existing[items[i].MenuID] = &items[i]
	}

	for _, ri := range reqItems {
		// Remove item if quantity = 0
		if ri.Quantity == 0 {
			if _, found := existing[ri.MenuID]; found {
				delete(existing, ri.MenuID)
			}
			continue
		}

		menu, err := u.MenuRepo.GetMenuByID(ctx, ri.MenuID)
		if err != nil {
			return nil, 0, fmt.Errorf("menu item '%s' not found: %w", ri.MenuID, err)
		}

		if menu.Status != "available" {
			return nil, 0, fmt.Errorf("menu item '%s' is not available", menu.Name)
		}

		if item, found := existing[ri.MenuID]; found {
			item.Quantity = ri.Quantity
			item.Status = domain.OrderItemStatusPending
			item.Notes = ri.Notes
			item.Subtotal = float64(item.Quantity) * item.Price
		} else {
			newItem := domain.OrderItem{
				MenuID:   menu.ID,
				MenuName: menu.Name,
				Price:    menu.Price,
				Quantity: ri.Quantity,
				Subtotal: float64(ri.Quantity) * menu.Price,
				Status:   domain.OrderItemStatusPending,
				Notes:    ri.Notes,
			}
			items = append(items, newItem)
			existing[ri.MenuID] = &items[len(items)-1]
		}
	}

	// Rebuild items slice from map to apply deletions
	items = items[:0]
	for _, item := range existing {
		items = append(items, *item)
	}

	for _, item := range items {
		total += item.Subtotal
	}

	return items, total, nil
}

func (u *OrderUsecase) UpdateOrderItemNote(ctx context.Context, req dto.UpdateOrderItemNoteRequest) (*domain.Order, error) {
	// 1. Fetch existing order
	existing, err := u.Repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// 2. Guard against updating a terminal order
	if existing.Status.IsTerminal() {
		return nil, fmt.Errorf("cannot update item note on order with status %q", existing.Status)
	}

	// 3. Persist
	result, err := u.Repo.UpdateOrderItemNote(ctx, req.OrderID, req.MenuID, req.Note)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item note: %w", err)
	}

	return result, nil
}

func (u *OrderUsecase) GetAllPendingOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	items, err := u.Repo.GetAllPendingOrderItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending order items: %w", err)
	}

	return items, nil
}

func (u *OrderUsecase) GetAllReadyOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	items, err := u.Repo.GetAllReadyOrderItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ready order items: %w", err)
	}

	return items, nil
}

func (u *OrderUsecase) GetActiveOrders(ctx context.Context) ([]domain.Order, error) {
	items, err := u.Repo.GetActiveOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ready order items: %w", err)
	}

	return items, nil
}

func (u *OrderUsecase) GetAllServedOrderItems(ctx context.Context) ([]domain.OrderItem, error) {
	items, err := u.Repo.GetAllServedOrderItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get served order items: %w", err)
	}

	return items, nil
}

func (u *OrderUsecase) UpdateOrderItemStatus(ctx context.Context, req dto.UpdateOrderItemStatusRequest) (*domain.Order, error) {
	// 1. Fetch existing order
	existing, err := u.Repo.GetOrderByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// 2. Guard against terminal order
	if existing.Status.IsTerminal() {
		return nil, fmt.Errorf("cannot update item status on order with status %q", existing.Status)
	}

	// 3. Validate status transition
	if !req.Status.IsValid() {
		return nil, fmt.Errorf("invalid order item status: %q", req.Status)
	}

	// 4. Persist
	result, err := u.Repo.UpdateOrderItemStatus(ctx, req.OrderID, req.MenuID, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update order item status: %w", err)
	}

	return result, nil
}

func generateOrderNumber(t time.Time) string {
	// e.g. ORD-20240315-XXXX (last 4 chars of a uuid for uniqueness)
	suffix := strings.ToUpper(uuid.NewString()[:4])
	return fmt.Sprintf("ORD-%s-%s", t.Format("20060102"), suffix)
}

func NowWIB() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

// func generateOrderNumber() string {
//     now := time.Now()
//     // Format: ORD-YYYYMMDD-Last4CharsOfUUID
//     datePart := now.Format("20060102")
//     uniquePart := strings.ToUpper(uuid.New().String()[:4])
//     return fmt.Sprintf("ORD-%s-%s", datePart, uniquePart)
// }
