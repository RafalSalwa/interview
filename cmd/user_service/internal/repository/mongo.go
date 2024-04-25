package repository

import (
	"context"
	"errors"

	"github.com/RafalSalwa/auth-api/pkg/models"
	apiMongo "github.com/RafalSalwa/auth-api/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type MongoAdapter struct {
	DB         *mongo.Client
	cfg        apiMongo.Config
	collection *mongo.Collection
}

func NewMongoDBUserRepository(db *mongo.Client, cfg apiMongo.Config) UserRepository {
	return &MongoAdapter{
		DB:         db,
		cfg:        cfg,
		collection: db.Database(cfg.Database).Collection("users"),
	}
}

func (m *MongoAdapter) Save(ctx context.Context, user *models.UserDBModel) error {
	return nil
}

func (m *MongoAdapter) Update(ctx context.Context, user *models.UserDBModel) error {
	ctx, span := otel.GetTracerProvider().Tracer("mongodb").Start(ctx, "Repository/UpdateLastLogin")
	defer span.End()

	var um models.UserMongoModel
	if err := um.FromDBModel(user); err != nil {
		return err
	}

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	if err := m.collection.FindOneAndUpdate(ctx, um, bson.M{"$set": user}, ops).Decode(&um); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (m *MongoAdapter) FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("mongodb").Start(ctx, "Repository/FindOne")
	defer span.End()

	um := &models.UserMongoModel{}
	if err := um.FromDBModel(user); err != nil {
		return nil, err
	}
	if err := m.collection.FindOne(ctx, &um).Decode(&um); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	if err := user.FromMongoUser(um); err != nil {
		return nil, err
	}
	return user, nil
}

func (m *MongoAdapter) SingUp(user *models.UserDBModel) error {
	return nil
}

func (m *MongoAdapter) Load(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	return nil, nil
}

func (m *MongoAdapter) ById(ctx context.Context, id int64) (*models.UserDBModel, error) {
	return nil, nil
}

func (m *MongoAdapter) ByLogin(ctx context.Context, user *models.SignInUserRequest) (*models.UserDBModel, error) {
	return nil, nil
}

func (m *MongoAdapter) ConfirmVerify(ctx context.Context, vCode string) error {
	return nil
}

func (m *MongoAdapter) UpdateLastLogin(ctx context.Context, u *models.UserDBModel) (*models.UserDBModel, error) {
	return nil, nil
}

func (m *MongoAdapter) FindUserByID(uid int64) (*models.UserDBModel, error) {
	return nil, nil
}

func (m *MongoAdapter) ChangePassword(ctx context.Context, userid int64, password string) error {
	return nil
}
