package services

import (
	"retrobytes/internal/repos"
)

type CartService struct {
	Carts *repos.CartRepo
	Prods *repos.ProductRepo
}

func NewCartService(carts *repos.CartRepo, prods *repos.ProductRepo) *CartService {
	return &CartService{Carts: carts, Prods: prods}
}

func (s *CartService) Add(sessionID, productID string, qty int) error {
	if qty < 1 {
		qty = 1
	}
	cartID, err := s.Carts.EnsureCart(sessionID)
	if err != nil {
		return err
	}
	p, err := s.Prods.Get(productID)
	if err != nil {
		return err
	}
	return s.Carts.UpsertItem(cartID, productID, qty, p.Price)
}

type CartView struct {
	Items []repos.CartItemRow
	Total float64
}

func (s *CartService) View(sessionID string) (CartView, error) {
	cartID, err := s.Carts.EnsureCart(sessionID)
	if err != nil {
		return CartView{}, err
	}
	items, total, err := s.Carts.View(cartID)
	if err != nil {
		return CartView{}, err
	}
	return CartView{Items: items, Total: total}, nil
}
