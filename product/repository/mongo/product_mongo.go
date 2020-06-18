package mongo

import (
	"context"

	"github.com/iqdf/benjerry-service/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const collectionName = "IceCream" // products

// ProductModel ...
type ProductModel struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty"`
	ProductID            string             `bson:"productId,omitempty"`
	Name                 string             `bson:"name,omitempty"`
	ImageClosedURL       string             `bson:"imageclosed_url,omitempty"`
	ImageOpenURL         string             `bson:"imageopen_url,omitempty"`
	Description          string             `bson:"description,omitempty"`
	Story                string             `bson:"story,omitempty"`
	SourcingValues       *[]string          `bson:"sourcing_values,omitempty"`
	Ingredients          *[]string          `bson:"ingredients,omitempty"`
	AllergyInfo          string             `bson:"allergy_info,omitempty"`
	DietaryCertification string             `bson:"dietary_certifications,omitempty"`
}

// ProductMongoRepo ...
type ProductMongoRepo struct {
	client *mongo.Client
	db     *mongo.Database
}

// modelFromProduct creates new ProductModel and
// copy product data from a data tranfer object into model's attributes
func modelFromProduct(product domain.Product) ProductModel {
	return ProductModel{
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
}

// Product yields data transfer object domain.Product
func (model *ProductModel) Product() domain.Product {
	return domain.Product{
		ProductID:            model.ProductID,
		Name:                 model.Name,
		ImageClosedURL:       model.ImageClosedURL,
		ImageOpenURL:         model.ImageOpenURL,
		Description:          model.Description,
		Story:                model.Story,
		SourcingValues:       model.SourcingValues,
		Ingredients:          model.Ingredients,
		AllergyInfo:          model.AllergyInfo,
		DietaryCertification: model.DietaryCertification,
	}
}

// NewProductRepo ...
func NewProductRepo(client *mongo.Client, dbName string) *ProductMongoRepo {
	repo := &ProductMongoRepo{
		client: client,
		db:     client.Database(dbName),
	}

	collection := repo.db.Collection(collectionName)

	// create unique index for product_id field
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "productId", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)
	return repo
}

// Fetch queries paginated products
// TODO: also return the next page?
// TODO: also input limit per page?
func (repo *ProductMongoRepo) Fetch(ctx context.Context) ([]domain.Product, error) {
	return []domain.Product{}, nil
}

// Get queries a single product identified by productID
func (repo *ProductMongoRepo) Get(ctx context.Context, productID string) (domain.Product, error) {
	var model ProductModel

	collection := repo.db.Collection(collectionName)
	err := collection.FindOne(ctx, ProductModel{ProductID: productID}).Decode(&model)

	return model.Product(), translate(err)
}

// Create inserts a single product document into collection
func (repo *ProductMongoRepo) Create(ctx context.Context, product domain.Product) error {
	var model = modelFromProduct(product)

	collection := repo.db.Collection(collectionName)
	_, err := collection.InsertOne(ctx, model)

	return translate(err)
}

// Update modifies attribute of a single product document
func (repo *ProductMongoRepo) Update(ctx context.Context, productID string, product domain.Product) error {
	var model = modelFromProduct(product)

	collection := repo.db.Collection(collectionName)
	filter := ProductModel{ProductID: productID}

	update := bson.M{"$set": model}
	_, err := collection.UpdateOne(ctx, filter, update)

	return translate(err)
}

// Delete removes a single document from collection
func (repo *ProductMongoRepo) Delete(ctx context.Context, productID string) error {
	collection := repo.db.Collection(collectionName)
	filter := ProductModel{ProductID: productID}

	_, err := collection.DeleteOne(ctx, filter)
	return translate(err)
}
