package config

import (
	"github.com/RafalSalwa/auth-api/pkg/csrf"
	"github.com/RafalSalwa/auth-api/pkg/env"
	"github.com/RafalSalwa/auth-api/pkg/http"
	"github.com/RafalSalwa/auth-api/pkg/http/auth"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/probes"
	"github.com/spf13/viper"
)

type (
	Config struct {
		ServiceName string         `mapstructure:"serviceName"`
		App         App            `mapstructure:"app"`
		Logger      *logger.Config `mapstructure:"logger"`
		HTTP        http.Config    `mapstructure:"http"`
		Auth        auth.Auth      `mapstructure:"auth"`
		Grpc        Grpc           `mapstructure:"grpc"`
		Probes      probes.Config  `mapstructure:"probes"`
		CSRF        csrf.Config    `mapstructure:"csrf"`
	}
	App struct {
		Env   string `mapstructure:"env"`
		Debug bool   `mapstructure:"debug"`
	}
	Grpc struct {
		AuthServicePort string `mapstructure:"authServicePort"`
		UserServicePort string `mapstructure:"userServicePort"`
	}
)

func InitConfig() (*Config, error) {
	path, err := env.GetConfigPath("gateway")
	if err != nil {
		return nil, err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
