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
	// 1. Resolve and validate each menu item
	items, total, err := u.resolveOrderItems(ctx, req.Items)
	if err != nil {
		return nil, err
	}

	// 2. Build the order domain object
	now := time.Now().UTC()
	order := &domain.Order{
		ID:          generateOrderID(), // e.g. "order::" + uuid
		Type:        "order",
		BookingID:   req.BookingID,
		OrderNumber: generateOrderNumber(now), // e.g. "ORD-20240315-0001"
		Items:       items,
		TotalAmount: total,
		Status:      "pending",
		CreatedAt:   NowWIB(),
	}

	// 3. Persist
	result, err := u.Repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return result, nil
}

func (u *OrderUsecase) resolveOrderItems(ctx context.Context, reqItems []dto.OrderItemRequest) ([]domain.OrderItem, float64, error) {
	var items []domain.OrderItem
	var total float64

	for _, ri := range reqItems {
		menu, err := u.MenuRepo.GetMenuByID(ctx, ri.MenuID)
		if err != nil {
			return nil, 0, fmt.Errorf("menu item '%s' not found: %w", ri.MenuID, err)
		}
		if menu.Status != "available" {
			return nil, 0, fmt.Errorf("menu item '%s' is not available", menu.Name)
		}

		subtotal := menu.Price * float64(ri.Quantity)
		total += subtotal

		items = append(items, domain.OrderItem{
			MenuID:   menu.ID,
			MenuName: menu.Name,
			Price:    menu.Price,
			Quantity: ri.Quantity,
			Subtotal: subtotal,
			Notes:    ri.Notes,
		})
	}

	return items, total, nil
}

func generateOrderID() string {
	return "order::" + uuid.NewString()
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
