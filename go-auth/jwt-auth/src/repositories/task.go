package repositories

type TaskRepository struct {
	*MongoRepositoryContext
}

func NewTaskRepository(uri, dbName, collectionName string) (*TaskRepository, error) {
	mongoRepo, err := NewMongoRepositoryContext(uri, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	return &TaskRepository{MongoRepositoryContext: mongoRepo}, nil
}
