package http

import (
	"encoding/json"
	"net/http"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"

	"github.com/iqdf/benjerry-service/domain"
	validatorLib "github.com/iqdf/benjerry-service/validator"
)

// productSingleResponse ...
type productSingleResponse struct {
	Data productData `json:"product"`
}

// productMultiResponse ...
type productMultiResponse struct {
	Data []productData `json:"products"`
}

type productData struct {
	ProductID            string    `json:"productId"`
	Name                 string    `json:"name"`
	ImageClosedURL       string    `json:"image_closed"`
	ImageOpenURL         string    `json:"image_open"`
	Description          string    `json:"description"`
	Story                string    `json:"story"`
	SourcingValues       *[]string `json:"sourcing_values"`
	Ingredients          *[]string `json:"ingredients"`
	AllergyInfo          string    `json:"allergy_info"`
	DietaryCertification string    `json:"dietary_certifications"`
}

// MessageError ....
type MessageError struct {
	Message string `json:"message"`
}

type productCreateRequest struct {
	ProductID            string    `json:"productId" validate:"required,numeric,min=3"`
	Name                 string    `json:"name" validate:"required,ascii,max=50"`
	ImageClosedURL       string    `json:"image_closed" validate:"omitempty,uri"`
	ImageOpenURL         string    `json:"image_open" validate:"omitempty,uri"`
	Description          string    `json:"description" validate:"ascii,max=100"`
	Story                string    `json:"story" validate:"omitempty,ascii,max=300"`
	SourcingValues       *[]string `json:"sourcing_values"`
	Ingredients          *[]string `json:"ingredients"`
	AllergyInfo          string    `json:"allergy_info" validate:"omitempty,ascii,max=50"`
	DietaryCertification string    `json:"dietary_certifications" validate:"omitempty,ascii,max=25"`
}

type productUpdateRequest struct {
	Name                 string    `json:"name" validate:"omitempty,ascii,max=50"`
	ImageClosedURL       string    `json:"image_closed" validate:"omitempty,uri"`
	ImageOpenURL         string    `json:"image_open" validate:"omitempty,uri"`
	Description          string    `json:"description" validate:"omitempty,ascii,max=100"`
	Story                string    `json:"story" validate:"omitempty,ascii,max=300"`
	SourcingValues       *[]string `json:"sourcing_values"`
	Ingredients          *[]string `json:"ingredients"`
	AllergyInfo          string    `json:"allergy_info" validate:"omitempty,ascii,max=50"`
	DietaryCertification string    `json:"dietary_certifications" validate:"omitempty,ascii,max=25"`
}

func createToProduct(requestData productCreateRequest) domain.Product {
	return domain.Product{
		ProductID:            requestData.ProductID,
		Name:                 requestData.Name,
		ImageClosedURL:       requestData.ImageClosedURL,
		ImageOpenURL:         requestData.ImageOpenURL,
		Description:          requestData.Description,
		Story:                requestData.Story,
		SourcingValues:       requestData.SourcingValues,
		Ingredients:          requestData.Ingredients,
		AllergyInfo:          requestData.AllergyInfo,
		DietaryCertification: requestData.DietaryCertification,
	}
}

func updateToProduct(requestData productUpdateRequest) domain.Product {
	return domain.Product{
		Name:                 requestData.Name,
		ImageClosedURL:       requestData.ImageClosedURL,
		ImageOpenURL:         requestData.ImageOpenURL,
		Description:          requestData.Description,
		Story:                requestData.Story,
		SourcingValues:       requestData.SourcingValues,
		Ingredients:          requestData.Ingredients,
		AllergyInfo:          requestData.AllergyInfo,
		DietaryCertification: requestData.DietaryCertification,
	}
}

// ProductHandler ...
type ProductHandler struct {
	service  domain.ProductService
	validate *validator.Validate
	trans    ut.Translator
}

// NewProductHandler creates new HTTP handler
// for product related request
func NewProductHandler(service domain.ProductService) *ProductHandler {
	validate, trans := validatorLib.NewValidator()
	handler := &ProductHandler{
		service:  service,
		validate: validate,
		trans:    trans,
	}
	return handler
}

func newSingleResponse(product domain.Product) productSingleResponse {
	productData := productData{
		ProductID:            product.ProductID,
		Name:                 product.Name,
		ImageClosedURL:       product.ImageClosedURL,
		ImageOpenURL:         product.ImageOpenURL,
		Description:          product.Description,
		Story:                product.Story,
		SourcingValues:       product.SourcingValues,
		Ingredients:          product.Ingredients,
		AllergyInfo:          product.AllergyInfo,
		DietaryCertification: product.DietaryCertification,
	}
	return productSingleResponse{Data: productData}
}

// Routes register handle func with the path url
func (handler *ProductHandler) Routes(router *mux.Router) {
	// Register handler methods to router here...
	router.HandleFunc("/{product_id}", handler.handleGetProduct()).Methods("GET")
	router.HandleFunc("/{product_id}", handler.handleUpdateProduct()).Methods("PUT")
	router.HandleFunc("/{product_id}", handler.handleDeleteProduct()).Methods("DELETE")
	router.HandleFunc("/", handler.handleCreateProduct()).Methods("POST")
}

// handleGetProduct provides handler func that gets a product
// [GET] /api/products/:product_id
func (handler *ProductHandler) handleGetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		params := mux.Vars(r)
		productID := params["product_id"]

		product, err := handler.service.GetProduct(r.Context(), productID)

		if err != nil {
			status := getAppErrorStatus(err)
			writeErrorMessage(w, err.Error(), status)
			return
		}

		response := newSingleResponse(product)
		json.NewEncoder(w).Encode(response)
	}
}

// handleCreateProduct provides handler func that creates a product
// [POST] /api/product/
func (handler *ProductHandler) handleCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var productCreate productCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&productCreate); err != nil {
			msg := getErrorMessage(err, handler.trans)
			writeErrorMessage(w, msg, http.StatusBadRequest)
			return
		}

		if errs := handler.validate.Struct(productCreate); errs != nil {
			msg := getErrorMessage(errs, handler.trans)
			writeErrorMessage(w, msg, http.StatusBadRequest)
			return
		}

		var product = createToProduct(productCreate)
		err := handler.service.CreateProduct(r.Context(), product)

		if err != nil {
			status := getAppErrorStatus(err)
			writeErrorMessage(w, err.Error(), status)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// handleUpdateProduct provides handler func that updates a product
// [PUT] /api/product/:product_id
func (handler *ProductHandler) handleUpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		productID := params["product_id"]

		var productUpdate productUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&productUpdate); err != nil {
			msg := getErrorMessage(err, handler.trans)
			writeErrorMessage(w, msg, http.StatusBadRequest)
			return
		}

		var product = updateToProduct(productUpdate)
		product.ProductID = productID
		err := handler.service.UpdateProduct(r.Context(), productID, product)

		if err != nil {
			status := getAppErrorStatus(err)
			msg := getErrorMessage(err, handler.trans)
			writeErrorMessage(w, msg, status)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// handleDeleteProduct provides handler func that deletes a product
// [DEL] /api/product/:product_id
func (handler *ProductHandler) handleDeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")

		params := mux.Vars(r)
		productID := params["product_id"]

		err := handler.service.DeleteProduct(r.Context(), productID)

		if err != nil {
			status := getAppErrorStatus(err)
			msg := getErrorMessage(err, handler.trans)
			writeErrorMessage(w, msg, status)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// getErrorMessage infers the error struct type
// and translates the error message to get client-friendly message
func getErrorMessage(err error, trans ut.Translator) string {
	if err == nil {
		return ""
	}

	switch err.(type) {
	case *json.UnmarshalTypeError:
		e, _ := err.(*json.UnmarshalTypeError)
		return e.Field + " field must be of type " + e.Type.String()

	case validator.ValidationErrors:
		var fieldErrs []string
		for _, e := range err.(validator.ValidationErrors) {
			fieldErrs = append(fieldErrs, e.Translate(trans))
		}
		return strings.Join(fieldErrs, "\n")
	}

	return err.Error()
}

// writerErrorMessage is a helper that writes error message to response
func writeErrorMessage(writer http.ResponseWriter, errMsg string, httpStatus int) {
	writer.WriteHeader(httpStatus)
	json.NewEncoder(writer).Encode(MessageError{Message: errMsg})
}

// getAppErrorStatus inputs error from application
// and infers the appropriate HTTP status to be returned
func getAppErrorStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var status int
	switch err.(type) {
	case *json.SyntaxError, *json.UnmarshalTypeError, *validator.ValidationErrors:
		return http.StatusBadRequest
	}

	if err == domain.ErrResourceNotFound { // TODO: check against error
		status = http.StatusNotFound
	} else if err == domain.ErrConflict {
		status = http.StatusConflict
	} else {
		status = http.StatusInternalServerError
	}

	return status
}
