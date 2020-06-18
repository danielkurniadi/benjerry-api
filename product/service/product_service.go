package service

import (
	"context"
	"time"

	"github.com/iqdf/benjerry-service/domain"
)

const timeout = time.Second * 10

// ProductService ...
type ProductService struct {
	productRepo domain.ProductRepository
}

// NewProductService creates new service
// that manages resource in product repository
func NewProductService(productRepo domain.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

// GetProduct ...
func (service *ProductService) GetProduct(ctx context.Context, productID string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	product, err := service.productRepo.Get(ctx, productID)

	if err != nil {
		return domain.Product{}, err
	}

	return product, nil
}

// CreateProduct ...
func (service *ProductService) CreateProduct(ctx context.Context, product domain.Product) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := service.productRepo.Create(ctx, product)

	if err != nil {
		return err
	}

	return nil
}

// UpdateProduct ...
func (service *ProductService) UpdateProduct(ctx context.Context, productID string, product domain.Product) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := service.productRepo.Update(ctx, productID, product)

	if err != nil {
		return err
	}

	return nil
}

// DeleteProduct ...
func (service *ProductService) DeleteProduct(ctx context.Context, productID string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	err := service.productRepo.Delete(ctx, productID)

	if err != nil {
		return err
	}

	return nil
}
