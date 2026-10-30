package cartservice

import (
	"sync"

	"github.com/verdex4/go-course-ozon/grpc/pkg/api/product"
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
	return cart
}

func (s *Storage) AddProduct(userID int64, product *product.Response, count uint32) Cart {
	cart := s.GetCart(userID)
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range cart.Items {
		if cart.Items[i].Sku == product.Sku {
			cart.Items[i].Count += count
			cart.TotalPrice += product.Price * count
			s.carts[userID] = cart
			return cart
		}
	}
	cart.Items = append(cart.Items, CartItem{
		Sku:   product.Sku,
		Name:  product.Name,
		Count: count,
		Price: product.Price,
	})

	cart.TotalPrice += product.Price * count
	s.carts[userID] = cart
	return cart
}
