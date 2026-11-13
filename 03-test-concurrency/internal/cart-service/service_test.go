package cartservice

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	product "github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockProductClient struct {
	product.ProductServiceClient
	getProductFunc func(ctx context.Context, in *product.Request, opts ...grpc.CallOption) (*product.Response, error)
}

func (m *mockProductClient) GetProduct(ctx context.Context, in *product.Request, opts ...grpc.CallOption) (*product.Response, error) {
	if m.getProductFunc != nil {
		return m.getProductFunc(ctx, in, opts...)
	}
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func setupRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{user_id}/cart", h.GetCartHandler)
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", h.AddProductHandler)
	return mux
}

func TestGetCartHandler(t *testing.T) {
	t.Parallel()

	t.Run("returns empty cart for non-existent user", func(t *testing.T) {
		t.Parallel()
		h := &Handler{productClient: &mockProductClient{}, storage: NewStorage()}
		router := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/user/1/cart", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var cart Cart
		err := json.Unmarshal(w.Body.Bytes(), &cart)
		require.NoError(t, err)
		require.Empty(t, cart.Items)
		require.Equal(t, uint32(0), cart.TotalPrice)
	})

	t.Run("returns existing cart with items", func(t *testing.T) {
		t.Parallel()

		mockClient := &mockProductClient{
			getProductFunc: func(ctx context.Context, in *product.Request, opts ...grpc.CallOption) (*product.Response, error) {
				return &product.Response{Sku: in.Sku, Name: "Test", Price: 50}, nil
			},
		}

		h := &Handler{productClient: mockClient, storage: NewStorage()}

		p := &product.Response{Sku: 100, Name: "Test", Price: 50}
		h.storage.AddProduct(1, p, 2)

		router := setupRouter(h)
		req := httptest.NewRequest(http.MethodGet, "/user/1/cart", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var cart Cart
		err := json.Unmarshal(w.Body.Bytes(), &cart)
		require.NoError(t, err)
		require.Len(t, cart.Items, 1)
		require.Equal(t, uint32(100), cart.TotalPrice) // 50 * 2
	})

	t.Run("bad request for invalid user_id", func(t *testing.T) {
		t.Parallel()
		h := &Handler{productClient: &mockProductClient{}, storage: NewStorage()}
		router := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/user/abc/cart", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddProductHandler(t *testing.T) {
	t.Parallel()

	t.Run("successfully adds product", func(t *testing.T) {
		t.Parallel()
		mockClient := &mockProductClient{
			getProductFunc: func(ctx context.Context, in *product.Request, opts ...grpc.CallOption) (*product.Response, error) {
				return &product.Response{Sku: in.Sku, Name: "Apple", Price: 10}, nil
			},
		}
		h := &Handler{productClient: mockClient, storage: NewStorage()}
		router := setupRouter(h)

		body := bytes.NewBufferString(`{"count": 5}`)
		req := httptest.NewRequest(http.MethodPost, "/user/1/cart/100", body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNoContent, w.Code)

		cart := h.storage.GetCart(1)
		require.Len(t, cart.Items, 1)
		require.Equal(t, uint32(5), cart.Items[0].Count)
	})

	t.Run("returns 412 when product not found", func(t *testing.T) {
		t.Parallel()
		mockClient := &mockProductClient{
			getProductFunc: func(ctx context.Context, in *product.Request, opts ...grpc.CallOption) (*product.Response, error) {
				return nil, status.Error(codes.NotFound, "product not found")
			},
		}
		h := &Handler{productClient: mockClient, storage: NewStorage()}
		router := setupRouter(h)

		body := bytes.NewBufferString(`{"count": 1}`)
		req := httptest.NewRequest(http.MethodPost, "/user/1/cart/999", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusPreconditionFailed, w.Code)
	})

	t.Run("bad request for invalid json body", func(t *testing.T) {
		t.Parallel()
		h := &Handler{productClient: &mockProductClient{}, storage: NewStorage()}
		router := setupRouter(h)

		body := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/user/1/cart/100", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad request for count <= 0", func(t *testing.T) {
		t.Parallel()
		h := &Handler{productClient: &mockProductClient{}, storage: NewStorage()}
		router := setupRouter(h)

		body := bytes.NewBufferString(`{"count": 0}`)
		req := httptest.NewRequest(http.MethodPost, "/user/1/cart/100", body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
