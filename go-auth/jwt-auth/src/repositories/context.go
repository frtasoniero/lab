package repositories

import (
	"context"
	"fmt"

	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos"
	dtoPages "github.com/frtasoniero/lab/go-auth/jwt-auth/src/dtos/pagination"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils"
	"github.com/frtasoniero/lab/go-auth/jwt-auth/src/utils/converts"
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

func (r *MongoRepositoryContext) GetPagination(ctx context.Context, params dtoPages.PaginationParams) dtoPages.PaginationResultContext {
	filter := bson.M{params.Field: params.Value}
	if params.SearchValue != "" {
		filter["$and"] = []bson.M{
			{params.Field: params.Value},
			{
				params.SearchField: bson.M{
					"$regex":   params.SearchValue,
					"$options": "i",
				},
			},
		}
	}

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return dtoPages.PaginationResultContext{
			TotalPages:  0,
			HasNextPage: false,
			TotalItems:  0,
			Err:         utils.InternalServerError(fmt.Sprintf("Error while counting documents: %v", err)),
		}
	}

	totalPages := int(count) / params.Limit
	if int(count)%params.Limit > 0 {
		totalPages++
	}

	if params.Skip < 1 {
		params.Skip = 1
	}

	if params.SortField == "" {
		params.SortField = "_id"
	}

	skipInit := (params.Skip - 1) * params.Limit

	hasNextPage := params.Skip < totalPages
	cursor, err := r.Collection.Find(
		ctx,
		filter,
		options.Find().SetSkip(int64(skipInit)).SetLimit(int64(params.Limit)).SetSort(bson.M{params.SortField: params.SortOrder}),
	)
	if err != nil {
		return dtoPages.PaginationResultContext{
			TotalPages:  0,
			HasNextPage: false,
			TotalItems:  skipInit,
			Err:         utils.InternalServerError(fmt.Sprintf("Error while searching for registers: %v", err)),
		}
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, params.Result)
	if err != nil {
		return dtoPages.PaginationResultContext{
			TotalPages:  0,
			HasNextPage: false,
			TotalItems:  skipInit,
			Err:         utils.InternalServerError(fmt.Sprintf("Error while decoding registers: %v", err)),
		}
	}

	return dtoPages.PaginationResultContext{
		TotalPages:  totalPages,
		HasNextPage: hasNextPage,
		TotalItems:  int(count),
		Err:         nil,
	}
}

func (r *MongoRepositoryContext) Update(contextServer context.Context, params dtos.UptadeFilter) error {
	filter := bson.M{"_id": params.Id}
	convertedToMap := converts.MapToKeyAndValueUptade(params.Dto)
	updateDoc := bson.M{"$set": convertedToMap}
	if params.ForeignKey != "" && params.ForeignKeyValue != primitive.NilObjectID {
		filter["$and"] = []bson.M{
			{"_id": params.Id},
			{params.ForeignKey: params.ForeignKeyValue},
		}
	}
	result, err := r.Collection.UpdateOne(contextServer, filter, updateDoc)
	if err != nil {
		return utils.BadRequestError(fmt.Sprintf("erro ao atualizar documento: %v", err))
	}
	if result.MatchedCount == 0 {
		return utils.NotFoundError("Nenhum documento encontrado para atualizar")
	}
	return nil
}
