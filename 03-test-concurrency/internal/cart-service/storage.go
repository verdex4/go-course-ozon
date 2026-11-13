package cartservice

import (
	"sync"

	"github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
)

type CartItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type Cart struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type Storage struct {
	mu    sync.RWMutex
	carts map[int64]Cart
}

func NewStorage() *Storage {
	return &Storage{
		carts: make(map[int64]Cart),
	}
}

func (s *Storage) GetCart(userID int64) Cart {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, exists := s.carts[userID]
	if !exists {
		return Cart{
			Items:      []CartItem{},
			TotalPrice: 0,
		}
	}

	itemsCopy := make([]CartItem, len(cart.Items))
	copy(itemsCopy, cart.Items)

	return Cart{
		Items:      itemsCopy,
		TotalPrice: cart.TotalPrice,
	}
}

func (s *Storage) AddProduct(userID int64, prod *product.Response, count uint32) Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		cart = Cart{Items: []CartItem{}, TotalPrice: 0}
	}

	found := false
	for i := range cart.Items {
		if cart.Items[i].Sku == prod.Sku {
			cart.Items[i].Count += count
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, CartItem{
			Sku:   prod.Sku,
			Name:  prod.Name,
			Count: count,
			Price: prod.Price,
		})
	}

	cart.TotalPrice += prod.Price * count
	s.carts[userID] = cart
	return cart
}
