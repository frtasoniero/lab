package repositories

import (
	"context"
	"fmt"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoRepositoryContext wraps a MongoDB collection and client instance.
// It provides methods for interacting with a MongoDB collection.
type MongoRepositoryContext struct {
	*mongo.Collection
	client *mongo.Client
}

// NewMongoRepositoryContext creates a new MongoRepositoryContext by connecting
// to the MongoDB server at the given URI, accessing the specified database and collection.
//
// Parameters:
//   - uri: MongoDB connection URI.
//   - dbName: Name of the database.
//   - collectionName: Name of the collection.
//
// Returns:
//   - A pointer to MongoRepositoryContext on successful connection.
//   - An error if the connection or ping fails.
func NewMongoRepositoryContext(uri, dbName, collectionName string) (*MongoRepositoryContext, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, utils.InternalServerError(fmt.Sprintf("Error while trying to connect to MongoDB: %v", err))
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, utils.InternalServerError(fmt.Sprintf("Error while trying to check MongoDB connectivity: %v", err))
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoRepositoryContext{
		Collection: collection,
		client:     client,
	}, nil
}

// Create inserts a new document into the MongoDB collection.
//
// Parameters:
//   - contextServer: A context for managing the request lifecycle.
//   - document: The document to be inserted (must be BSON-marshalable).
//
// Returns:
//   - An error if the insertion fails, wrapped as a bad request error.
func (r *MongoRepositoryContext) Create(contextServer context.Context, document interface{}) error {
	_, err := r.Collection.InsertOne(contextServer, document)
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("Error while trying to insert document: %v", err))
	}

	return nil
}

func (r *MongoRepositoryContext) ExistsByAny(ctx context.Context, params dtos.ExistsFilter) error {
	filter := bson.M{params.Field: params.Value}
	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {
		filter["$and"] = []bson.M{
			{params.Field: params.Value},
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}

	err := r.Collection.FindOne(ctx, filter).Err()
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return utils.InternalServerError(fmt.Sprintf("Error while searching for the user: %v", err))
	}

	return utils.ConflictError(fmt.Sprintf("User already exists: %v", params.Value))
}

func (r *MongoRepositoryContext) GetByAny(ctx context.Context, params dtos.GetAnyFilter) error {
	filter := bson.M{params.Field: params.Value}

	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {
		filter["$and"] = []bson.M{
			{params.Field: params.Value},
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}

	err := r.Collection.FindOne(ctx, filter).Decode(params.Result)
	if err == mongo.ErrNoDocuments {
		return utils.NotFoundError("Register not found")
	} else if err != nil {
		return utils.InternalServerError(fmt.Sprintf("Error while searching for register: %v", err))
	}
	return nil
}
