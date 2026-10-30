package cartservice

import (
	"testing"

	"github.com/stretchr/testify/require"
	product "github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
)

func TestStorage_GetCart(t *testing.T) {
	t.Parallel()

	t.Run("returns empty cart for non-existent user", func(t *testing.T) {
		t.Parallel()

		s := NewStorage()

		cart := s.GetCart(1)

		require.Empty(t, cart.Items, "items should be empty")
		require.Equal(t, uint32(0), cart.TotalPrice, "total price should be 0")
	})

	t.Run("returns existing cart with items", func(t *testing.T) {
		t.Parallel()

		s := NewStorage()

		p := &product.Response{Sku: 100, Name: "Test Product", Price: 50}
		s.AddProduct(1, p, 2)

		cart := s.GetCart(1)

		require.Len(t, cart.Items, 1)
		require.Equal(t, int64(100), cart.Items[0].Sku)
		require.Equal(t, "Test Product", cart.Items[0].Name)
		require.Equal(t, uint32(2), cart.Items[0].Count)
		require.Equal(t, uint32(100), cart.TotalPrice) // 50 * 2
	})
}

func TestStorage_AddProduct(t *testing.T) {
	t.Parallel()

	t.Run("adds new product to non-existent user", func(t *testing.T) {
		t.Parallel()

		s := NewStorage()
		p := &product.Response{Sku: 1, Name: "Apple", Price: 10}

		cart := s.AddProduct(1, p, 5)

		require.Len(t, cart.Items, 1)
		require.Equal(t, int64(1), cart.Items[0].Sku)
		require.Equal(t, uint32(5), cart.Items[0].Count)
		require.Equal(t, uint32(50), cart.TotalPrice) // 10 * 5
	})

	t.Run("increments count of existing product", func(t *testing.T) {
		t.Parallel()

		s := NewStorage()
		p := &product.Response{Sku: 1, Name: "Apple", Price: 10}

		s.AddProduct(1, p, 2)
		cart := s.AddProduct(1, p, 3)

		require.Len(t, cart.Items, 1, "should not create duplicate items")
		require.Equal(t, uint32(5), cart.Items[0].Count) // 2 + 3
		require.Equal(t, uint32(50), cart.TotalPrice)    // 10 * 5
	})

	t.Run("adds different products to same user", func(t *testing.T) {
		t.Parallel()

		s := NewStorage()
		p1 := &product.Response{Sku: 1, Name: "Apple", Price: 10}
		p2 := &product.Response{Sku: 2, Name: "Banana", Price: 20}

		s.AddProduct(1, p1, 1)
		cart := s.AddProduct(1, p2, 2)

		require.Len(t, cart.Items, 2)
		// 10*1 + 20*2 = 50
		require.Equal(t, uint32(50), cart.TotalPrice)
	})
}
