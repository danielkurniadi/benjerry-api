package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/iqdf/benjerry-service/common/auth"
	mongoHelper "github.com/iqdf/benjerry-service/common/repository/mongo"
	"github.com/iqdf/benjerry-service/domain"
)

const collectionName = "User"

// UserModel ...
type UserModel struct {
	ID            primitive.ObjectID    `bson:"_id,omitempty"`
	Username      string                `bson:"username,omitempty"`
	HashPassword  string                `bson:"hashpassword,omitempty"`
	Authorization *[]auth.Authorization `bson:"authorizations,omitempty"`
}

// UserMongoRepo ...
type UserMongoRepo struct {
	client *mongo.Client
	db     *mongo.Database
}

// modelFromUser creates new UserModel and
// copy data from user entity to user DB model
func modelFromUser(user domain.User) UserModel {
	return UserModel{
		Username:      user.Username,
		HashPassword:  user.HashPassword,
		Authorization: &user.Authorizations,
	}
}

// User creates user entity instance and
// copies data from model into user entity
func (model *UserModel) User() domain.User {
	return domain.User{
		Username:       model.Username,
		HashPassword:   model.HashPassword,
		Authorizations: *model.Authorization,
	}
}

// NewUserRepo ...
func NewUserRepo(client *mongo.Client, dbName string) *UserMongoRepo {
	repo := &UserMongoRepo{
		client: client,
		db:     client.Database(dbName),
	}

	collection := repo.db.Collection(collectionName)

	// create unique index constraint for username field
	collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{Key: "username", Value: bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	return repo
}

// Get queries a single user identified by username
func (repo *UserMongoRepo) Get(ctx context.Context, username string) (domain.User, error) {
	var model UserModel

	collection := repo.db.Collection(collectionName)
	err := collection.FindOne(ctx, UserModel{Username: username}).Decode(&model)
	return model.User(), mongoHelper.TranslateError(err)
}

// Create inserts a single user document into collection
func (repo *UserMongoRepo) Create(ctx context.Context, user domain.User) error {
	var model UserModel = modelFromUser(user)

	collection := repo.db.Collection(collectionName)
	_, err := collection.InsertOne(ctx, model)

	return mongoHelper.TranslateError(err)
}
