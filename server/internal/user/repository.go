package user

import "go.uber.org/zap"

type Repository interface {
	FindById(id UserId) (*User, error)
	FindByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
}

type repositoryMiddleware func(Repository) Repository

type mongoRepository struct{}

func NewMongoRepository(logger *zap.Logger) Repository {
	return &mongoRepository{}
}

func (mongorepository *mongoRepository) FindById(id UserId) (*User, error) {
	panic("not implemented") // TODO: Implement
}

func (mongorepository *mongoRepository) FindByUsername(username string) (*User, error) {
	panic("not implemented") // TODO: Implement
}

func (mongorepository *mongoRepository) Create(user *User) error {
	panic("not implemented") // TODO: Implement
}

func (mongorepository *mongoRepository) Update(user *User) error {
	panic("not implemented") // TODO: Implement
}
