package cmd

import (
	"log"
	"net"

	productservice "github.com/verdex4/go-course-ozon/test-concurrency/internal/product-service"
	"github.com/verdex4/go-course-ozon/test-concurrency/pkg/api/product"
	"google.golang.org/grpc"
)

var (
	jsonFilePath = "internal/product-service/data.json"
	port         = ":50051"
)

func StartProductService() {
	storage, err := productservice.NewStorage(jsonFilePath)
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}
	log.Printf("Successfully loaded products from %s", jsonFilePath)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	srvc := productservice.NewService(storage)
	product.RegisterProductServiceServer(grpcServer, srvc)

	log.Printf("Product service gRPC running on port %s...", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
