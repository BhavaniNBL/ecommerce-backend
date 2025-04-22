package service

import (
	"context"
	"log"

	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/repository"
)

type ProductService struct {
	Repo *repository.ProductRepository
	productpb.UnimplementedProductServiceServer
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) Create(product *model.Product) error {
	return s.Repo.Create(product)
}

func (s *ProductService) GetByID(id string) (*model.Product, error) {
	// return s.Repo.GetByID(id)
	product, err := s.Repo.GetByID(id)
	if err != nil {
		log.Println("GetByID error:", err) // <-- Add this
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Update(product *model.Product) error {
	return s.Repo.Update(product)
}

func (s *ProductService) Delete(id string) error {
	return s.Repo.Delete(id)
}

func (s *ProductService) List(filters map[string]string) ([]model.Product, error) {
	return s.Repo.List(filters)
}

// In product-service.go
func (s *ProductService) ListProducts(ctx context.Context, _ *productpb.Empty) (*productpb.ProductList, error) {
	products, err := s.Repo.List(nil)
	if err != nil {
		return nil, err
	}

	var grpcProducts []*productpb.Product
	for _, p := range products {
		grpcProducts = append(grpcProducts, p.ToGRPC())
	}

	return &productpb.ProductList{Products: grpcProducts}, nil
}
