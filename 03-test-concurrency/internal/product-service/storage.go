package productservice

import (
	"encoding/json"
	"os"
	"strconv"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type ProductModel struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Storage struct {
	products map[int64]ProductModel
}

func NewStorage(filePath string) (*Storage, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rawProducts := make(map[string]ProductModel)
	if err := json.Unmarshal(fileData, &rawProducts); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	productsMap := make(map[int64]ProductModel, len(rawProducts))
	for sku, p := range rawProducts {
		skuInt, err := strconv.ParseInt(sku, 10, 64)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		productsMap[skuInt] = p
	}

	return &Storage{
		products: productsMap,
	}, nil
}

func (s *Storage) Get(sku int64) (ProductModel, bool) {
	product, exists := s.products[sku]
	return product, exists
}
