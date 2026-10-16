package productservice

import (
	context "context"

	product "github.com/verdex4/go-course-ozon/grpc/pkg/api/product"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Service struct {
	product.UnimplementedProductServiceServer
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetProduct(ctx context.Context, req *product.Request) (*product.Response, error) {
	reqSku := req.GetSku()
	if reqSku <= 0 {
		return nil, status.Error(codes.InvalidArgument, "sku must be greater than zero")
	}

	p, exists := s.storage.Get(req.GetSku())
	if !exists {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &product.Response{
		Sku:   reqSku,
		Name:  p.Name,
		Price: p.Price,
	}, nil
}
