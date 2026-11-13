//go:build api

package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	product "github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestProductService_Integration(t *testing.T) {
	client, cleanup, err := NewProductClient()
	require.NoError(t, err)
	t.Cleanup(cleanup)

	t.Run("GetProduct returns existing product", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := NewContext()
		defer cancel()

		// SKU 1076963 есть в data.json
		resp, err := client.GetProduct(ctx, &product.Request{Sku: 1076963})

		require.NoError(t, err)
		require.Equal(t, int64(1076963), resp.Sku)
		require.NotEmpty(t, resp.Name)
		require.Greater(t, resp.Price, uint32(0))
	})

	t.Run("GetProduct returns NotFound for non-existent SKU", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := NewContext()
		defer cancel()

		nonExistentSKU := int64(999999999)
		resp, err := client.GetProduct(ctx, &product.Request{Sku: nonExistentSKU})

		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.NotFound, status.Code(err))
	})

	t.Run("GetProduct returns InvalidArgument for SKU <= 0", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := NewContext()
		defer cancel()

		resp, err := client.GetProduct(ctx, &product.Request{Sku: 0})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))

		resp, err = client.GetProduct(ctx, &product.Request{Sku: -1})
		require.Error(t, err)
		require.Nil(t, resp)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	})
}
