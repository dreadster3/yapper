package steps

import (
	"context"
	"errors"

	"github.com/dreadster3/yapper/server/internal/domain"
	"github.com/dreadster3/yapper/server/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type StepRepository interface {
	GetByMessageId(ctx context.Context, messageId domain.MessageId) ([]*Step, error)
	Create(ctx context.Context, step *Step) error
}

const collectionName = "steps"

type step struct {
	Id        primitive.ObjectID `bson:"_id"`
	MessageId primitive.ObjectID `bson:"message_id"`
	Type      string             `bson:"type"`
	Content   string             `bson:"content"`
	Status    string             `bson:"status"`
}

func (s step) ToModel() *Step {
	return &Step{
		Id:        domain.StepId(s.Id.Hex()),
		MessageId: domain.MessageId(s.MessageId.Hex()),
		Type:      s.Type,
		Content:   s.Content,
		Status:    StepStatus(s.Status),
	}
}

func fromModel(m *Step) (*step, error) {
	id, err := primitive.ObjectIDFromHex(string(m.Id))
	if err != nil {
		if errors.Is(err, primitive.ErrInvalidHex) {
			id = primitive.NewObjectID()
		}
	}

	messageId, err := primitive.ObjectIDFromHex(string(m.MessageId))
	if err != nil {
		return nil, err
	}

	return &step{
		Id:        id,
		MessageId: messageId,
		Type:      m.Type,
		Content:   m.Content,
		Status:    string(m.Status),
	}, nil
}

type stepRepository struct {
	db *mongo.Database
}

func NewStepRepository(db *mongo.Database, logger *zap.Logger) StepRepository {
	var repo StepRepository
	repo = &stepRepository{db: db}
	repo = NewLoggingMiddleware(logger)(repo)
	return repo
}

func (r *stepRepository) Collection() *mongo.Collection {
	return r.db.Collection(collectionName)
}

func (r *stepRepository) GetByMessageId(ctx context.Context, messageId domain.MessageId) ([]*Step, error) {
	objId, err := primitive.ObjectIDFromHex(string(messageId))
	if err != nil {
		return nil, err
	}

	cursor, err := r.Collection().Find(ctx, bson.M{"message_id": objId})
	if err != nil {
		return nil, err
	}

	var steps []step
	if err := cursor.All(ctx, &steps); err != nil {
		return nil, err
	}

	return utils.Map(steps, func(s step) *Step {
		return s.ToModel()
	}), nil
}

func (r *stepRepository) Create(ctx context.Context, step *Step) error {
	entity, err := fromModel(step)
	if err != nil {
		return err
	}

	result, err := r.Collection().InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	step.Id = domain.StepId(result.InsertedID.(primitive.ObjectID).Hex())
	return nil
}

type repositoryMiddleware func(StepRepository) StepRepository

type loggingMiddleware struct {
	logger *zap.Logger
	next   StepRepository
}

func NewLoggingMiddleware(logger *zap.Logger) repositoryMiddleware {
	return func(next StepRepository) StepRepository {
		return &loggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

func (m *loggingMiddleware) GetByMessageId(ctx context.Context, messageId domain.MessageId) (steps []*Step, err error) {
	defer func() {
		m.logger.Debug("GetByMessageId", zap.String("message_id", string(messageId)), zap.Objects("steps", steps), zap.Error(err))
	}()

	return m.next.GetByMessageId(ctx, messageId)
}

func (m *loggingMiddleware) Create(ctx context.Context, step *Step) (err error) {
	defer func() {
		m.logger.Debug("Create", zap.Object("step", step), zap.Error(err))
	}()

	return m.next.Create(ctx, step)
}
