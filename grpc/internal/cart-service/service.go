package cartservice

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	product "github.com/verdex4/go-course-ozon/grpc/pkg/api/product"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
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
	if err != nil {
		http.Error(w, "user_id must be an integer", http.StatusBadRequest)
		return
	}
	if userID <= 0 {
		http.Error(w, "user_id must be greater than zero", http.StatusBadRequest)
		return
	}

	cart := h.storage.GetCart(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

// POST /user/{user_id}/cart/{sku_id}
// Body: {"count": <count>}
func (h *Handler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "user_id must be an integer", http.StatusBadRequest)
		return
	}
	if userID <= 0 {
		http.Error(w, "user_id must be greater than zero", http.StatusBadRequest)
		return
	}
	skuIDStr := r.PathValue("sku_id")
	skuID, err := strconv.ParseInt(skuIDStr, 10, 64)
	if err != nil {
		http.Error(w, "sku_id must be an integer", http.StatusBadRequest)
		return
	}
	if skuID <= 0 {
		http.Error(w, "sku_id must be greater than zero", http.StatusBadRequest)
		return
	}

	type bodyRequest struct {
		Count uint32 `json:"count"`
	}
	var req bodyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Count <= 0 {
		http.Error(w, "count must be greater than zero", http.StatusBadRequest)
		return
	}

	product, err := h.productClient.GetProduct(r.Context(), &product.Request{Sku: skuID})
	if errors.Is(err, status.Error(codes.NotFound, "product not found")) {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.storage.AddProduct(userID, product, req.Count)
	w.WriteHeader(http.StatusNoContent)
}
