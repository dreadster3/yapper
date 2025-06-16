package chats

import (
	"context"
	"errors"

	"github.com/dreadster3/yapper/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *Chat) error
	FindById(ctx context.Context, id ChatId) (*Chat, error)
}

const (
	collectionName = "chats"
)

type chat struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	ProfileId primitive.ObjectID `bson:"profile_id"`
}

func (c chat) ToModel() *Chat {
	return &Chat{
		Id:        ChatId(c.Id.Hex()),
		Name:      c.Name,
		ProfileId: domain.ProfileId(c.ProfileId.Hex()),
	}
}

func fromModel(c *Chat) (*chat, error) {
	id, err := primitive.ObjectIDFromHex(string(c.Id))
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			id = primitive.NewObjectID()
		}
	}

	profileId, err := primitive.ObjectIDFromHex(string(c.ProfileId))
	if err != nil {
		return nil, err
	}

	return &chat{
		Id:        id,
		Name:      c.Name,
		ProfileId: profileId,
	}, nil
}

type chatRepository struct {
	db *mongo.Database
}

func NewChatRepository(db *mongo.Database, logger *zap.Logger) ChatRepository {
	var repo ChatRepository
	repo = &chatRepository{db: db}
	repo = NewLoggingMiddleware(logger)(repo)
	return repo
}

func (r *chatRepository) Collection() *mongo.Collection {
	return r.db.Collection(collectionName)
}

func (r *chatRepository) Create(ctx context.Context, chat *Chat) error {
	entity, err := fromModel(chat)
	if err != nil {
		return err
	}

	result, err := r.Collection().InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	chat.Id = ChatId(result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

func (r *chatRepository) FindById(ctx context.Context, id ChatId) (*Chat, error) {
	objId, err := primitive.ObjectIDFromHex(string(id))
	if err != nil {
		return nil, err
	}

	var chat *chat
	if err := r.Collection().FindOne(ctx, bson.M{"_id": objId}).Decode(&chat); err != nil {
		return nil, err
	}

	return chat.ToModel(), nil
}

type repositoryMiddleware func(ChatRepository) ChatRepository

type loggingMiddleware struct {
	logger *zap.Logger
	next   ChatRepository
}

func NewLoggingMiddleware(logger *zap.Logger) repositoryMiddleware {
	return func(next ChatRepository) ChatRepository {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (m *loggingMiddleware) Create(ctx context.Context, chat *Chat) (err error) {
	defer func() {
		m.logger.Debug("Create", zap.Object("chat", chat), zap.Error(err))
	}()

	return m.next.Create(ctx, chat)
}

func (m *loggingMiddleware) FindById(ctx context.Context, id ChatId) (chat *Chat, err error) {
	defer func() {
		m.logger.Debug("FindById", zap.String("id", string(id)), zap.Object("chat", chat), zap.Error(err))
	}()

	return m.next.FindById(ctx, id)
}
