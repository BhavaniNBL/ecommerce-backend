package service

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"gorm.io/gorm"

	"github.com/BhavaniNBL/ecommerce-backend/proto/inventorypb"
	"github.com/BhavaniNBL/ecommerce-backend/proto/productpb"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/model"
	"github.com/BhavaniNBL/ecommerce-backend/services/product-service/repository"
)

type ProductService struct {
	Repo            *repository.ProductRepository
	InventoryClient inventorypb.InventoryServiceClient
	productpb.UnimplementedProductServiceServer
}

//	func NewProductService(repo *repository.ProductRepository) *ProductService {
//		return &ProductService{Repo: repo}
//	}
func NewProductService(repo *repository.ProductRepository, invClient inventorypb.InventoryServiceClient) *ProductService {
	return &ProductService{
		Repo:            repo,
		InventoryClient: invClient,
	}
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

// product_service.go

// func (s *ProductService) CheckProductExists(ctx context.Context, req *productpb.ProductID) (*productpb.ProductExistsResponse, error) {
// 	// Use the repository to check if the product exists
// 	_, err := s.Repo.GetByID(req.Id)
// 	if err != nil {
// 		// If the error is "record not found", product doesn't exist
// 		// if err.Error() == "record not found" {
// 		// 	return &productpb.ProductExistsResponse{Exists: false}, nil
// 		// }
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return &productpb.ProductExistsResponse{Exists: false}, nil
// 		}
// 		// Log other errors
// 		log.Println("Error checking product existence:", err)
// 		return nil, err
// 	}

// 	// If we got here, the product exists
// 	return &productpb.ProductExistsResponse{Exists: true}, nil
// }

func (s *ProductService) CheckProductExists(ctx context.Context, req *productpb.ProductID) (*productpb.ProductExistsResponse, error) {
	log.Println("📦 [CheckProductExists] Called with ID:", req.Id)

	product, err := s.Repo.GetByID(req.Id)
	if err != nil {
		log.Println("❌ [CheckProductExists] GetByID error:", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &productpb.ProductExistsResponse{Exists: false}, nil
		}
		return nil, err
	}

	log.Println("✅ [CheckProductExists] Product found:", product.ID)
	return &productpb.ProductExistsResponse{Exists: true}, nil
}

// CheckInventoryForProduct calls the Inventory Service to check the stock for a product
func (s *ProductService) CheckInventoryForProduct(ctx context.Context, productID string, token string) (*inventorypb.InventoryResponse, error) {
	// Create metadata with Authorization header
	// md := metadata.New(map[string]string{
	// 	"authorization": "Bearer " + token, // Pass the token here
	// })

	// Create a new context with metadata
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)

	// Create a gRPC connection to the Inventory service
	conn, err := grpc.NewClient("inventory-service:50053", grpc.WithTransportCredentials(insecure.NewCredentials())) // Use the appropriate address
	if err != nil {
		log.Fatalf("Failed to connect to Inventory service: %v", err)
		return nil, err
	}
	defer conn.Close()

	// Create a client for the Inventory service
	client := inventorypb.NewInventoryServiceClient(conn)

	// Call the Inventory service (assuming you have an InventoryRequest)
	resp, err := client.GetInventory(ctx, &inventorypb.GetInventoryRequest{
		ProductId: productID,
	})
	if err != nil {
		log.Fatalf("Error calling Inventory service: %v", err)
		return nil, err
	}

	return resp, nil
}
