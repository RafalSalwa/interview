package repository

import (
	"context"

	"github.com/RafalSalwa/auth-api/cmd/user_service/config"
	"github.com/RafalSalwa/auth-api/pkg/mongo"
	"github.com/RafalSalwa/auth-api/pkg/sql"
)

const (
	MySQL   string = "mysql"
	MongoDB string = "mongo"
)

func NewUserRepository(ctx context.Context, dbType string, params *config.Config) (UserRepository, error) {
	switch dbType {
	case MySQL:
		con, err := sql.NewGormConnection(params.MySQL)
		if err != nil {
			return nil, err
		}
		return NewUserAdapter(con), nil

	case MongoDB:
		mongoClient, err := mongo.NewClient(ctx, params.Mongo)
		if err != nil {
			return nil, err
		}
		return newMongoDBUserRepository(mongoClient, params.Mongo), nil

	default:
		panic("Unsupported database type")
	}
}
