package service

import (
	"context"
	"testing"

	"github.com/iqdf/benjerry-service/domain"
	"github.com/iqdf/benjerry-service/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	contextType   = mock.Anything
	productType   = mock.AnythingOfType("domain.Product")
	productIDType = mock.AnythingOfType("string")
)

func TestGetByProductID(t *testing.T) {
	// setup mock repository and mock item
	mockProductRepo := new(mocks.ProductRepository)

	t.Run("GetProduct-success", func(t *testing.T) {
		mockProductSuccess := createMockProduct()
		mockProductRepo.On("Get", contextType, productIDType).
			Return(mockProductSuccess, nil).
			Once()

		var productService = NewProductService(mockProductRepo)
		product, err := productService.GetProduct(context.TODO(), mockProductSuccess.ProductID)

		assert.NoError(t, err)
		assert.Equal(t, product.Name, mockProductSuccess.Name) // TODO: use cmp.Equal
	})

	t.Run("GetProduct-on-db-error", func(t *testing.T) {
		var dberr error = domain.ErrResourceNotFound

		mockProductID := "646"
		mockProductFail := domain.Product{}
		mockProductRepo.On("Get", contextType, productIDType).
			Return(mockProductFail, dberr).
			Once()

		var productService = NewProductService(mockProductRepo)
		product, err := productService.GetProduct(context.TODO(), mockProductID)

		assert.Error(t, err)
		assert.Equal(t, dberr, err)
		assert.Equal(t, product, mockProductFail) // TODO: use cmp.Equal
	})
}

func TestCreateProduct(t *testing.T) {
	// setup mock repository and mock data
	mockProductRepo := new(mocks.ProductRepository)

	t.Run("CreateProduct-success", func(t *testing.T) {
		mockProductSuccess := createMockProduct()
		mockProductRepo.On("Create", contextType, productType).
			Return(nil).
			Once()

		var productService = NewProductService(mockProductRepo)
		err := productService.CreateProduct(context.TODO(), mockProductSuccess)

		assert.NoError(t, err)
	})

	t.Run("CreateProduct-on-db-error", func(t *testing.T) {
		var dberr error = domain.ErrConflict
		mockProductFail := domain.Product{}

		mockProductRepo.On("Create", contextType, productType).
			Return(dberr).
			Once()

		var productService = NewProductService(mockProductRepo)
		err := productService.CreateProduct(context.TODO(), mockProductFail)

		assert.Error(t, err)
		assert.Equal(t, dberr, err)
	})
}

func TestUpdateProduct(t *testing.T) {
	// setup mock repository and mock data
	mockProductRepo := new(mocks.ProductRepository)

	t.Run("UpdateProduct-success", func(t *testing.T) {
		mockProductSuccess := createMockProduct()
		mockProductID := mockProductSuccess.ProductID

		mockProductRepo.On("Update", contextType, productIDType, productType).
			Return(nil).
			Once()

		var productService = NewProductService(mockProductRepo)
		err := productService.UpdateProduct(context.TODO(), mockProductID, mockProductSuccess)

		assert.NoError(t, err)
	})

	t.Run("UpdateProduct-on-db-error", func(t *testing.T) {
		var dberr error = domain.ErrResourceNotFound
		mockProductFail := domain.Product{}

		mockProductRepo.On("Update", contextType, productIDType, productType).
			Return(dberr).
			Once()

		var productService = NewProductService(mockProductRepo)
		err := productService.UpdateProduct(context.TODO(), "646", mockProductFail)

		assert.Error(t, err)
		assert.Equal(t, dberr, err)
	})
}

func TestDeleteProduct(t *testing.T) {
	// setup mock repository and mock data
	mockProductRepo := new(mocks.ProductRepository)

	t.Run("DeleteProduct-success", func(t *testing.T) {
		mockProductID := "646"
		mockProductRepo.On("Delete", contextType, productIDType).
			Return(nil).
			Once()

		var productService = NewProductService(mockProductRepo)
		err := productService.DeleteProduct(context.TODO(), mockProductID)

		assert.NoError(t, err)
	})

	t.Run("DeleteProduct-on-db-error", func(t *testing.T) {
		var dberr error = domain.ErrResourceNotFound
		mockProductID := "646"

		mockProductRepo.On("Delete", contextType, productIDType).
			Return(dberr).
			Once()

		var productService = NewProductService(mockProductRepo)
		err := productService.DeleteProduct(context.TODO(), mockProductID)

		assert.Error(t, err)
		assert.Equal(t, dberr, err)
	})
}

func createMockProduct() domain.Product {
	mockProductSuccess := domain.Product{
		ProductID:      "646",
		Name:           "Vanilla Toffee Bar Crunch",
		ImageClosedURL: "/files/vanilla-toffee-landing.png",
		ImageOpenURL:   "/files/vanilla-toffee-landing-open.png",
		Description:    "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
		Story:          "Vanilla What Bar Crunch? The new ice cream with toffee bars...",
		SourcingValues: &[]string{
			"Non-GMO",
			"Cage-Free Eggs",
			"Fairtrade",
			"Responsibly Sourced Packaging",
			"Caring Dairy",
		},
		Ingredients: &[]string{
			"cream",
			"skim milk",
			"liquid sugar",
			"water",
			"sugar",
		},
		AllergyInfo:          "May contain wheat, peanuts",
		DietaryCertification: "Singapore Food Ministry",
	}
	return mockProductSuccess
}
