package db

import (
	"context"
	"errors"
	"go-grpc-gateway/pb"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type GormProvider struct {
	db *gorm.DB
}

// NewGormProvider initializes the GormProvider.
func NewGormProvider(db *gorm.DB) *GormProvider {
	return &GormProvider{db: db}
}

// AddProduct creates a new product in the database.
func (p *GormProvider) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.AddProductResponse, error) {
	now := time.Now()
	product := &pb.ProductORM{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		Category:    req.Category,
		Stock:       req.Stock,
		Image:       req.Image,
		CreatedAt:   &now, // Assign the address of `now` to `*time.Time` fields
		UpdatedAt:   nil,
	}
	log.Println(product)

	if err := p.db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, err
	}
	p.db = p.db.Debug()

	return &pb.AddProductResponse{
		ProductId: strconv.FormatUint(uint64(product.Id), 10),
		Message:   "Product added successfully",
	}, nil
}

// GetProduct fetches a product by its ID.
func (p *GormProvider) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	var product pb.ProductORM
	if err := p.db.WithContext(ctx).First(&product, "id = ?", req.ProductId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &pb.GetProductResponse{
		ProductId:   strconv.FormatUint(uint64(product.Id), 10),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		Stock:       product.Stock,
		Currency:    product.Currency,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt.String(),
		UpdatedAt:   product.UpdatedAt.String(),
	}, nil
}

// UpdateProduct updates the details of a product.
func (p *GormProvider) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	var product pb.ProductORM
	if err := p.db.WithContext(ctx).First(&product, "id = ?", req.ProductId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	now := time.Now()
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.UpdatedAt = &now

	if err := p.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, err
	}

	return &pb.UpdateProductResponse{Message: "Product updated successfully"}, nil
}

// DeleteProduct removes a product by its ID.
func (p *GormProvider) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if err := p.db.WithContext(ctx).Delete(&pb.ProductORM{}, "id = ?", req.ProductId).Error; err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{Message: "Product deleted successfully"}, nil
}

// ListProduct retrieves all products.
func (p *GormProvider) ListProduct(ctx context.Context, req *pb.ListProductRequest) (*pb.ListProductResponse, error) {
	var products []pb.ProductORM
	if err := p.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}

	var productResponses []*pb.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, &pb.GetProductResponse{
			ProductId:   strconv.FormatUint(uint64(product.Id), 10),
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Category:    product.Category,
			Stock:       product.Stock,
			Image:       product.Image,
			CreatedAt:   product.CreatedAt.String(),
			UpdatedAt:   product.UpdatedAt.String(),
		})
	}

	return &pb.ListProductResponse{Products: productResponses}, nil
}
