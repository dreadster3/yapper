package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/dreadster3/yapper/server/internal/platform/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var (
	ErrProfileNotFound       = errors.New("profile not found")
	ErrProfileAlreadyCreated = errors.New("profile already created")
)

type ProfileRepository interface {
	FindById(ctx context.Context, id ProfileId) (*Profile, error)
	FindByUserId(ctx context.Context, userId auth.UserId) (*Profile, error)
	Create(ctx context.Context, profile *Profile) error
}

const (
	collection = "profiles"
)

type profileEntity struct {
	Id     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	UserId string             `bson:"user_id"`
}

func (p profileEntity) ToModel() *Profile {
	return &Profile{
		Id:     ProfileId(p.Id.Hex()),
		Name:   p.Name,
		UserId: auth.UserId(p.UserId),
	}
}

func fromModel(p *Profile) *profileEntity {
	id, err := primitive.ObjectIDFromHex(string(p.Id))
	if err != nil && errors.Is(err, primitive.ErrInvalidHex) {
		id = primitive.NewObjectID()
	}

	return &profileEntity{
		Id:     id,
		Name:   p.Name,
		UserId: string(p.UserId),
	}
}

type profileRepository struct {
	db *mongo.Database
}

func NewProfileRepository(db *mongo.Database, logger *zap.Logger) ProfileRepository {
	var repo ProfileRepository

	repo = &profileRepository{db: db}
	repo = NewLoggingMiddleware(logger)(repo)

	return repo
}

func (r *profileRepository) Collection() *mongo.Collection {
	return r.db.Collection(collection)
}

func (r *profileRepository) FindById(ctx context.Context, id ProfileId) (*Profile, error) {
	return nil, errors.New("not implemented") // TODO: Implement
}

func (r *profileRepository) FindByUserId(ctx context.Context, userId auth.UserId) (*Profile, error) {
	var entity profileEntity
	if err := r.Collection().FindOne(ctx, bson.M{"user_id": userId}).Decode(&entity); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrProfileNotFound
		}

		return nil, fmt.Errorf("repository.FindByUserId: %w", err)
	}

	return entity.ToModel(), nil
}

func (r *profileRepository) Create(ctx context.Context, profile *Profile) error {
	if _, err := r.FindByUserId(ctx, profile.UserId); err != nil {
		if !errors.Is(err, ErrProfileNotFound) {
			return err
		}

		entity := fromModel(profile)
		result, err := r.Collection().InsertOne(ctx, entity)
		if err != nil {
			return fmt.Errorf("repository.Create: %w", err)
		}

		profile.Id = ProfileId(result.InsertedID.(primitive.ObjectID).Hex())
		return nil
	}

	return ErrProfileAlreadyCreated
}

type repositoryMiddleware func(ProfileRepository) ProfileRepository

type loggingMiddleware struct {
	logger *zap.Logger
	next   ProfileRepository
}

func NewLoggingMiddleware(logger *zap.Logger) repositoryMiddleware {
	return func(next ProfileRepository) ProfileRepository {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (m *loggingMiddleware) FindById(ctx context.Context, id ProfileId) (profile *Profile, err error) {
	defer func() {
		m.logger.Debug("FindById", zap.String("id", string(id)), zap.Object("profile", profile), zap.Error(err))
	}()

	return m.next.FindById(ctx, id)
}

func (m *loggingMiddleware) FindByUserId(ctx context.Context, userId auth.UserId) (profile *Profile, err error) {
	defer func() {
		m.logger.Debug("FindByUserId", zap.String("user_id", string(userId)), zap.Object("profile", profile), zap.Error(err))
	}()

	return m.next.FindByUserId(ctx, userId)
}

func (m *loggingMiddleware) Create(ctx context.Context, profile *Profile) error {
	defer func() {
		m.logger.Debug("Create", zap.Object("profile", profile))
	}()

	return m.next.Create(ctx, profile)
}
