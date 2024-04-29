package config

import (
	"errors"

	"github.com/RafalSalwa/auth-api/pkg/email"
	"github.com/RafalSalwa/auth-api/pkg/env"
	"github.com/RafalSalwa/auth-api/pkg/grpc"
	"github.com/RafalSalwa/auth-api/pkg/jwt"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	mongodb "github.com/RafalSalwa/auth-api/pkg/mongo"
	"github.com/RafalSalwa/auth-api/pkg/probes"
	"github.com/RafalSalwa/auth-api/pkg/rabbitmq"
	"github.com/RafalSalwa/auth-api/pkg/redis"
	"github.com/RafalSalwa/auth-api/pkg/sql"
	"github.com/spf13/viper"
)

type (
	Config struct {
		ServiceName      string              `mapstructure:"serviceName"`
		App              App                 `mapstructure:"app"`
		Logger           *logger.Config      `mapstructure:"logger"`
		GRPC             grpc.Config         `mapstructure:"grpc"`
		Email            email.Config        `mapstructure:"email"`
		JWTToken         jwt.JWTConfig       `mapstructure:"jwt"`
		MySQL            sql.MySQL           `mapstructure:"mysql"`
		Mongo            mongodb.Config      `mapstructure:"mongo"`
		MongoCollections mongodb.Collections `mapstructure:"mongoCollections"`
		Redis            *redis.Config       `mapstructure:"redis"`
		Rabbit           rabbitmq.Config     `mapstructure:"rabbitmq"`
		Probes           probes.Config       `mapstructure:"probes"`
	}
	App struct {
		Env            string `mapstructure:"env"`
		Debug          bool   `mapstructure:"debug"`
		RepositoryType string `mapstructure:"repository_type"`
	}
)

func InitConfig() (*Config, error) {
	cfg := &Config{}
	path, err := env.GetConfigPath("auth_service")
	if err != nil {
		return nil, err
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.New("viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.New("viper.Unmarshal")
	}
	return cfg, nil
}
