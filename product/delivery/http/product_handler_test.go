package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/iqdf/benjerry-service/domain"
	"github.com/iqdf/benjerry-service/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	contextType   = mock.Anything
	productIDType = mock.AnythingOfType("string")
	productType   = mock.AnythingOfType("domain.Product")
)

func TestGetProductSuccess(t *testing.T) {
	productService := new(mocks.ProductService)
	mockProduct := createMockProduct()
	productID := mockProduct.ProductID

	productService.On("GetProduct", contextType, productIDType).
		Return(mockProduct, nil).
		Once()

	request, _ := http.NewRequest("GET", "/"+productID, strings.NewReader(""))
	vars := map[string]string{
		"product_id": productID,
	}
	request = mux.SetURLVars(request, vars)
	recorder := httptest.NewRecorder()

	productHandler := NewProductHandler(productService)
	getHandle := productHandler.handleGetProduct()

	var productResponse productSingleResponse

	getHandle(recorder, request)
	err := json.NewDecoder(recorder.Body).Decode(&productResponse)

	assert.NoError(t, err)
	assert.Equal(t, recorder.Code, 200)
	assert.Equal(t, productResponse, newSingleResponse(mockProduct))
}

func TestGetProductNotFound(t *testing.T) {
	productService := new(mocks.ProductService)
	productID := "978"

	productService.On("GetProduct", contextType, productIDType).
		Return(domain.Product{}, domain.ErrResourceNotFound).
		Once()

	request, _ := http.NewRequest("GET", "/api/products/"+productID, strings.NewReader(""))
	vars := map[string]string{
		"product_id": productID,
	}
	request = mux.SetURLVars(request, vars)
	recorder := httptest.NewRecorder()

	productHandler := NewProductHandler(productService)
	getHandle := productHandler.handleGetProduct()

	getHandle(recorder, request)
	assert.Equal(t, recorder.Code, 404)
}

func TestCreateProductSuccess(t *testing.T) {
	productService := new(mocks.ProductService)

	productService.On("CreateProduct", contextType, productType).
		Return(nil).
		Once()

	createReq := createMockCreateRequest()
	productbyte, err := json.Marshal(createReq)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/products/", strings.NewReader(string(productbyte)))
	recorder := httptest.NewRecorder()

	assert.NoError(t, err)
	productHandler := NewProductHandler(productService)
	createHandle := productHandler.handleCreateProduct()

	createHandle(recorder, request)
	assert.Equal(t, recorder.Code, 201)
}

func TestCreateProductConflict(t *testing.T) {
	productService := new(mocks.ProductService)

	productService.On("CreateProduct", contextType, productType).
		Return(domain.ErrConflict).
		Once()

	createReq := createMockCreateRequest()
	productbyte, err := json.Marshal(createReq)
	assert.NoError(t, err)

	request, err := http.NewRequest("POST", "/api/products/", strings.NewReader(string(productbyte)))
	recorder := httptest.NewRecorder()

	assert.NoError(t, err)
	productHandler := NewProductHandler(productService)
	createHandle := productHandler.handleCreateProduct()

	createHandle(recorder, request)
	assert.Equal(t, recorder.Code, 200)

	var msgErr messageError
	json.NewDecoder(recorder.Body).Decode(&msgErr)
	assert.Equal(t, msgErr.Message, domain.ErrConflict.Error())
}

func TestUpdateSuccess(t *testing.T) {
	productService := new(mocks.ProductService)

	productService.On("UpdateProduct", contextType, productIDType, productType).
		Return(nil).
		Once()

	updateReq := createMockUpdateRequest()
	productbyte, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	request, err := http.NewRequest("PUT", "/api/products/", strings.NewReader(string(productbyte)))
	recorder := httptest.NewRecorder()

	assert.NoError(t, err)
	productHandler := NewProductHandler(productService)
	updateHandle := productHandler.handleUpdateProduct()

	updateHandle(recorder, request)
	assert.Equal(t, recorder.Code, 200)
}

func TestUpdateFail(t *testing.T) {
	productService := new(mocks.ProductService)

	productService.On("UpdateProduct", contextType, productIDType, productType).
		Return(domain.ErrResourceNotFound).
		Once()

	updateReq := createMockUpdateRequest()
	productbyte, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	request, err := http.NewRequest("PUT", "/api/products/", strings.NewReader(string(productbyte)))
	recorder := httptest.NewRecorder()

	assert.NoError(t, err)
	productHandler := NewProductHandler(productService)
	updateHandle := productHandler.handleUpdateProduct()

	updateHandle(recorder, request)
	assert.Equal(t, recorder.Code, 404)

	var msgErr messageError
	json.NewDecoder(recorder.Body).Decode(&msgErr)
	assert.Equal(t, msgErr.Message, domain.ErrResourceNotFound.Error())
}

func TestDeleteSuccess(t *testing.T) {
	productService := new(mocks.ProductService)

	productService.On("DeleteProduct", contextType, productIDType).
		Return(nil).
		Once()

	updateReq := createMockUpdateRequest()
	productbyte, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	request, err := http.NewRequest("DELETE", "/api/products/", strings.NewReader(string(productbyte)))
	recorder := httptest.NewRecorder()

	assert.NoError(t, err)
	productHandler := NewProductHandler(productService)
	deleteHandle := productHandler.handleDeleteProduct()

	deleteHandle(recorder, request)
	assert.Equal(t, recorder.Code, 200)
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

func createMockCreateRequest() productCreateRequest {
	return productCreateRequest{
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
}

func createMockUpdateRequest() productUpdateRequest {
	return productUpdateRequest{
		Name:        "Vanilla Toffee Bar Crunch",
		Description: "Updated: Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
		Story:       "Updated: Vanilla What Bar Crunch? The new ice cream with toffee bars...",
		SourcingValues: &[]string{
			"Non-GMO",
			"Cage-Free Eggs",
		},
		Ingredients: &[]string{},
	}
}
