package repositories

type UserRepository struct {
	*MongoRepositoryContext
}

func NewUserRepository(uri, dbName, collectionName string) (*UserRepository, error) {
	mongoRepo, err := NewMongoRepositoryContext(uri, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	return &UserRepository{MongoRepositoryContext: mongoRepo}, nil
}
