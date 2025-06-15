package messages

import (
	"context"
	"errors"
	"time"

	"github.com/dreadster3/yapper/server/internal/chats"
	"github.com/dreadster3/yapper/server/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MessageRepository interface {
	FindById(ctx context.Context, id domain.MessageId) (*Message, error)
	GetByChatId(ctx context.Context, chatId chats.ChatId) ([]*Message, error)
	Create(ctx context.Context, message *Message) error
	Update(ctx context.Context, message *Message) error
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
	Status    string             `bson:"status"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

func (m message) ToModel() *Message {
	return &Message{
		Id:        domain.MessageId(m.Id.Hex()),
		ChatId:    chats.ChatId(m.ChatId.Hex()),
		Provider:  m.Provider,
		Model:     m.Model,
		Role:      MessageRole(m.Role),
		Content:   m.Content,
		Status:    MessageStatus(m.Status),
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

	status := string(m.Status)
	if m.Status == "" {
		status = string(MessageStatusPending)
	}

	return &message{
		Id:        id,
		ChatId:    chatId,
		Provider:  m.Provider,
		Model:     m.Model,
		Role:      string(m.Role),
		Content:   m.Content,
		Status:    status,
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

func (r *messageRepository) FindById(ctx context.Context, id domain.MessageId) (*Message, error) {
	return nil, errors.New("not implemented") // TODO: Implement
}

func (r *messageRepository) GetByChatId(ctx context.Context, chatId chats.ChatId) ([]*Message, error) {
	objId, err := primitive.ObjectIDFromHex(string(chatId))
	if err != nil {
		return nil, err
	}

	cursor, err := r.Collection().Find(ctx, bson.M{"chat_id": objId})
	if err != nil {
		return nil, err
	}

	var entities []message
	if err := cursor.All(ctx, &entities); err != nil {
		return nil, err
	}

	result := make([]*Message, len(entities))
	for idx, e := range entities {
		result[idx] = e.ToModel()
	}

	return result, nil
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

	message.Id = domain.MessageId(result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

func (r *messageRepository) Update(ctx context.Context, message *Message) error {
	entity, err := fromModel(message)
	if err != nil {
		return err
	}

	if _, err := r.Collection().UpdateByID(ctx, entity.Id, bson.M{"$set": entity}); err != nil {
		return err
	}

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

func (m *loggerMiddleware) FindById(ctx context.Context, id domain.MessageId) (message *Message, err error) {
	defer func() {
		m.logger.Debug("FindById", zap.String("id", string(id)), zap.Object("message", message), zap.Error(err))
	}()

	return m.next.FindById(ctx, id)
}

func (m *loggerMiddleware) GetByChatId(ctx context.Context, chatId chats.ChatId) (messages []*Message, err error) {
	defer func() {
		m.logger.Debug("GetByChatId", zap.String("chat_id", string(chatId)), zap.Objects("messages", messages), zap.Error(err))
	}()

	return m.next.GetByChatId(ctx, chatId)
}

func (m *loggerMiddleware) Create(ctx context.Context, message *Message) (err error) {
	defer func() {
		m.logger.Debug("Create", zap.Object("message", message), zap.Error(err))
	}()

	return m.next.Create(ctx, message)
}

func (m *loggerMiddleware) Update(ctx context.Context, message *Message) (err error) {
	defer func() {
		m.logger.Debug("Update", zap.Object("message", message), zap.Error(err))
	}()

	return m.next.Update(ctx, message)
}
