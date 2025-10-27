package services

import (
	"database/sql"
	"errors"
	"fmt"

	"retrobytes/internal/repos"

	"github.com/google/uuid"
)

type Contact struct {
	Name  string
	Email string
}

type OrderService struct {
	Carts  *repos.CartRepo
	Inv    *repos.InventoryRepo
	Orders *repos.OrderRepo
}

func NewOrderService(carts *repos.CartRepo, inv *repos.InventoryRepo, orders *repos.OrderRepo) *OrderService {
	return &OrderService{Carts: carts, Inv: inv, Orders: orders}
}

func (s *OrderService) Place(sessionID, region, fulfillment string, contact Contact) (string, error) {
	if region == "" {
		return "", errors.New("missing region")
	}
	if fulfillment == "" {
		fulfillment = "delivery"
	}

	cartID, err := s.Carts.EnsureCart(sessionID)
	if err != nil {
		return "", err
	}

	items, err := s.Carts.Items(cartID)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "", errors.New("cart empty")
	}

	// pre-check stock
	for _, it := range items {
		qty, err := s.Inv.Qty(it.ProductID, region)
		if err != nil && err != sql.ErrNoRows {
			return "", err
		}
		if qty < it.Qty {
			return "", fmt.Errorf("insufficient stock for %s (need %d, have %d)", it.ProductID, it.Qty, qty)
		}
	}

	// decrement
	for _, it := range items {
		if err := s.Inv.Decrement(it.ProductID, region, it.Qty); err != nil {
			return "", err
		}
	}

	// totals
	total := 0.0
	for _, it := range items {
		total += it.Price * float64(it.Qty)
	}

	// create order
	orderID := uuid.NewString()
	if err := s.Orders.Create(orderID, sessionID, region, fulfillment, contact.Name, contact.Email, total); err != nil {
		return "", err
	}
	for _, it := range items {
		if err := s.Orders.InsertItem(orderID, it.ProductID, it.Qty, it.Price, it.Condition); err != nil {
			return "", err
		}
	}
	_ = s.Carts.Clear(cartID)
	return orderID, nil
}
