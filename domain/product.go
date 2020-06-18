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
	GetProduct(context.Context, string) (Product, error)
	CreateProduct(context.Context, Product) error
	UpdateProduct(context.Context, string, Product) error
	DeleteProduct(context.Context, string) error
}

// ProductRepository ...
type ProductRepository interface {
	Create(context.Context, Product) error
	Get(context.Context, string) (Product, error)
	Update(context.Context, string, Product) error
	Delete(context.Context, string) error
}
