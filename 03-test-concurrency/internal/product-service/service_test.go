package productservice

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	product "github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestService_GetProduct(t *testing.T) {
	t.Parallel()

	t.Run("returns existing product", func(t *testing.T) {
		t.Parallel()

		storage := &Storage{
			products: map[int64]ProductModel{
				100: {Name: "Apple", Price: 10},
			},
		}
		service := NewService(storage)

		resp, err := service.GetProduct(context.Background(), &product.Request{Sku: 100})
		require.NoError(t, err)
		require.Equal(t, int64(100), resp.Sku)
		require.Equal(t, "Apple", resp.Name)
		require.Equal(t, uint32(10), resp.Price)
	})

	t.Run("returns NotFound for non-existent product", func(t *testing.T) {
		t.Parallel()

		storage := &Storage{
			products: map[int64]ProductModel{},
		}
		service := NewService(storage)

		resp, err := service.GetProduct(context.Background(), &product.Request{Sku: 999})
		require.Error(t, err)
		require.Nil(t, resp)

		require.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("returns InvalidArgument for sku <= 0", func(t *testing.T) {
		t.Parallel()

		storage := &Storage{
			products: map[int64]ProductModel{},
		}
		service := NewService(storage)

		resp, err := service.GetProduct(context.Background(), &product.Request{Sku: 0})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))

		resp, err = service.GetProduct(context.Background(), &product.Request{Sku: -1})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
