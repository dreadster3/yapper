package messages

import (
	"context"
	"errors"
	"time"

	"github.com/dreadster3/yapper/server/internal/chat"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MessageRepository interface {
	FindById(ctx context.Context, id MessageId) (*Message, error)
	GetByChatId(ctx context.Context, chatId chat.ChatId) ([]*Message, error)
	Create(ctx context.Context, message *Message) error
}

const (
	collectionName = "messages"
)

type message struct {
	Id        primitive.ObjectID `bson:"_id"`
	ChatId    primitive.ObjectID `bson:"chat_id"`
	Provider  string             `bson:"provider"`
	Model     string             `bson:"model"`
	Role      string             `bson:"role"`
	Content   string             `bson:"content"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

func (m message) ToModel() *Message {
	return &Message{
		Id:        MessageId(m.Id.Hex()),
		ChatId:    chat.ChatId(m.ChatId.Hex()),
		Provider:  m.Provider,
		Model:     m.Model,
		Role:      m.Role,
		Content:   m.Content,
		CreatedAt: m.CreatedAt.Time(),
	}
}

func fromModel(m *Message) (*message, error) {
	id, err := primitive.ObjectIDFromHex(string(m.Id))
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			id = primitive.NewObjectID()
		}
	}

	chatId, err := primitive.ObjectIDFromHex(string(m.ChatId))
	if err != nil {
		return nil, err
	}

	createdAt := m.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	return &message{
		Id:        id,
		ChatId:    chatId,
		Provider:  m.Provider,
		Model:     m.Model,
		Role:      m.Role,
		Content:   m.Content,
		CreatedAt: primitive.NewDateTimeFromTime(createdAt),
	}, nil
}

type messageRepository struct {
	db *mongo.Database
}

func NewMessageRepository(db *mongo.Database, logger *zap.Logger) MessageRepository {
	var repo MessageRepository
	repo = &messageRepository{db: db}
	repo = NewLoggingMiddleware(logger)(repo)
	return repo
}

func (r *messageRepository) Collection() *mongo.Collection {
	return r.db.Collection(collectionName)
}

func (r *messageRepository) FindById(ctx context.Context, id MessageId) (*Message, error) {
	return nil, errors.New("not implemented") // TODO: Implement
}

func (r *messageRepository) GetByChatId(ctx context.Context, chatId chat.ChatId) ([]*Message, error) {
	return nil, errors.New("not implemented") // TODO: Implement
}

func (r *messageRepository) Create(ctx context.Context, message *Message) error {
	entity, err := fromModel(message)
	if err != nil {
		return err
	}

	result, err := r.Collection().InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	message.Id = MessageId(result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

type repositoryMiddleware func(MessageRepository) MessageRepository

type loggerMiddleware struct {
	next   MessageRepository
	logger *zap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) repositoryMiddleware {
	return func(next MessageRepository) MessageRepository {
		return &loggerMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (m *loggerMiddleware) FindById(ctx context.Context, id MessageId) (message *Message, err error) {
	defer func() {
		m.logger.Debug("FindById", zap.String("id", string(id)), zap.Object("message", message), zap.Error(err))
	}()

	return m.next.FindById(ctx, id)
}

func (m *loggerMiddleware) GetByChatId(ctx context.Context, chatId chat.ChatId) (messages []*Message, err error) {
	defer func() {
		m.logger.Debug("GetByChatId", zap.String("chat_id", string(chatId)), zap.Objects("messages", messages), zap.Error(err))
	}()

	return m.next.GetByChatId(ctx, chatId)
}

func (m *loggerMiddleware) Create(ctx context.Context, message *Message) error {
	defer func() {
		m.logger.Debug("Create", zap.Object("message", message))
	}()

	return m.next.Create(ctx, message)
}
