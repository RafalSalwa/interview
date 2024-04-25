package repository

import (
	"context"
	"errors"
	"time"

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

func newMongoDBUserRepository(db *mongo.Client, cfg apiMongo.Config) UserRepository {
	return &MongoAdapter{
		DB:         db,
		cfg:        cfg,
		collection: db.Database(cfg.Database).Collection("users"),
	}
}

func (m MongoAdapter) Exists(ctx context.Context, udb *models.UserDBModel) bool {
	ctx, span := otel.GetTracerProvider().Tracer("mongodb").Start(ctx, "Repository/Exists")
	defer span.End()

	var um models.UserMongoModel
	if err := um.FromDBModel(udb); err != nil {
		return false
	}
	if err := m.collection.FindOne(ctx, um).Decode(&um); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		return false
	}

	return true
}

func (m MongoAdapter) Update(ctx context.Context, user models.UserDBModel) error {
	// TODO mongo implement me
	panic("mongo Update implement me")
}

func (m MongoAdapter) Save(ctx context.Context, user *models.UserDBModel) error {
	ctx, span := otel.GetTracerProvider().Tracer("mongodb repository").Start(ctx, "Service SignUpUser")
	defer span.End()

	mu := models.UserMongoModel{}
	err := mu.FromDBModel(user)
	if err != nil {
		return err
	}
	_, err = m.collection.InsertOne(ctx, mu, &options.InsertOneOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (m MongoAdapter) FindAll(ctx context.Context, user *models.UserDBModel) ([]models.UserDBModel, error) {
	// TODO mongo implement me
	panic("mongo FindAll implement me")
}

func (m MongoAdapter) FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("mongodb").Start(ctx, "Repository/FindOne")
	defer span.End()

	var um models.UserMongoModel
	if err := um.FromDBModel(user); err != nil {
		return nil, err
	}
	if err := m.collection.FindOne(ctx, um).Decode(&um); err != nil {
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

func (m MongoAdapter) GetOrCreate(ctx context.Context, id int64) (*models.UserDBModel, error) {
	// TODO mongo implement me
	panic("mongo GetOrCreate implement me")
}

func (m MongoAdapter) Confirm(ctx context.Context, udb *models.UserDBModel) error {
	// TODO mongo implement me
	panic("mongo Confirm implement me")
}

func (m MongoAdapter) UpdateLastLogin(ctx context.Context, user *models.UserDBModel) error {
	ctx, span := otel.GetTracerProvider().Tracer("mongodb").Start(ctx, "Repository/UpdateLastLogin")
	defer span.End()

	var um models.UserMongoModel
	if err := um.FromDBModel(user); err != nil {
		return err
	}
	now := time.Now()
	um.UpdatedAt = &now

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	if err := m.collection.FindOneAndUpdate(ctx, um, bson.M{"$set": user}, ops).Decode(&um); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	if err := user.FromMongoUser(um); err != nil {
		return err
	}

	return nil
}
