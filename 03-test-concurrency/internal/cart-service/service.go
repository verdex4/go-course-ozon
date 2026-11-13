package cartservice

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	product "github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Handler struct {
	productClient product.ProductServiceClient
	storage       *Storage
}

func NewHandler(productServiceAddr string) (*Handler, error) {
	conn, err := grpc.NewClient(productServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := product.NewProductServiceClient(conn)
	h := &Handler{
		productClient: client,
		storage:       NewStorage(),
	}
	return h, nil
}

// GET /user/{user_id}/cart
func (h *Handler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "user_id must be a positive integer", http.StatusBadRequest)
		return
	}

	cart := h.storage.GetCart(userID)

	if len(cart.Items) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cart)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	updatedItems := make([]CartItem, len(cart.Items))

	for i, item := range cart.Items {
		sku := item.Sku

		g.Go(func() error {
			prod, err := h.productClient.GetProduct(ctx, &product.Request{Sku: sku})
			if err != nil {
				return err
			}

			updatedItems[i] = CartItem{
				Sku:   prod.Sku,
				Name:  prod.Name,
				Count: item.Count,
				Price: prod.Price,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "product not found in catalog", http.StatusPreconditionFailed)
			return
		}
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "timeout fetching products", http.StatusGatewayTimeout)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var totalPrice uint32
	for _, item := range updatedItems {
		totalPrice += item.Price * item.Count
	}

	responseCart := Cart{
		Items:      updatedItems,
		TotalPrice: totalPrice,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseCart)
}

// POST /user/{user_id}/cart/{sku_id}
func (h *Handler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "user_id must be a positive integer", http.StatusBadRequest)
		return
	}

	skuIDStr := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil || skuID <= 0 {
		http.Error(w, "sku_id must be a positive integer", http.StatusBadRequest)
		return
	}

	type bodyRequest struct {
		Count uint32 `json:"count"`
	}
	var req bodyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Count <= 0 {
		http.Error(w, "invalid request body or count <= 0", http.StatusBadRequest)
		return
	}

	prod, err := h.productClient.GetProduct(r.Context(), &product.Request{Sku: skuID})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			http.Error(w, "product not found", http.StatusPreconditionFailed)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.storage.AddProduct(userID, prod, req.Count)
	w.WriteHeader(http.StatusNoContent)
}
