package cmd

import (
	"log"
	"net/http"

	cart "github.com/verdex4/go-course-ozon/test-concurrency/internal/cart-service"
)

func StartCartService() {
	// Создание обработчика с подключением к product-сервису
	addr := "localhost:50051"
	handler, err := cart.NewHandler(addr)
	if err != nil {
		log.Fatalf("failed to init cart handler: %v", err)
	}

	// Создание HTTP-сервера на порту 8080
	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{user_id}/cart", handler.GetCartHandler)
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handler.AddProductHandler)

	port := ":8080"
	log.Printf("Cart service HTTP running on port %s...", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
