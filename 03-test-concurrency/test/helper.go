//go:build api

package test

import (
	"context"
	"time"

	product "github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	productServiceAddr = "localhost:50051"
	cartServiceAddr    = "http://localhost:8080"
)

func NewProductClient() (product.ProductServiceClient, func(), error) {
	conn, err := grpc.NewClient(
		productServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	client := product.NewProductServiceClient(conn)
	cleanup := func() {
		conn.Close()
	}

	return client, cleanup, nil
}

func NewContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
