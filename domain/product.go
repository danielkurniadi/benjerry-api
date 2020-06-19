package domain

import (
	"context"
)

// Product domain
type Product struct {
	ProductID            string
	Name                 string
	ImageClosedURL       string
	ImageOpenURL         string
	Description          string
	Story                string
	SourcingValues       *[]string
	Ingredients          *[]string
	AllergyInfo          string
	DietaryCertification string
}

// ProductService ...
type ProductService interface {
	GetProduct(ctx context.Context, productID string) (Product, error)
	CreateProduct(ctx context.Context, product Product) error
	UpdateProduct(ctx context.Context, productID string, product Product) error
	DeleteProduct(ctx context.Context, productID string) error
}

// ProductRepository ...
type ProductRepository interface {
	Create(ctx context.Context, product Product) error
	Get(ctx context.Context, productID string) (Product, error)
	Update(ctx context.Context, productID string, product Product) error
	Delete(ctx context.Context, productID string) error
}
