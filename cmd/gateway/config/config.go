package config

import (
	"github.com/RafalSalwa/interview-app-srv/pkg/csrf"
	"github.com/RafalSalwa/interview-app-srv/pkg/env"
	"github.com/RafalSalwa/interview-app-srv/pkg/http"
	"github.com/RafalSalwa/interview-app-srv/pkg/http/auth"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/probes"
	"github.com/RafalSalwa/interview-app-srv/pkg/tracing"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string               `mapstructure:"serviceName"`
	App         App                  `mapstructure:"app"`
	Logger      *logger.Config       `mapstructure:"logger"`
	HTTP        http.Config          `mapstructure:"http"`
	Auth        auth.Auth            `mapstructure:"auth"`
	Grpc        Grpc                 `mapstructure:"grpc"`
	Probes      probes.Config        `mapstructure:"probes"`
	Jaeger      tracing.JaegerConfig `mapstructure:"jaeger"`
	CSRF        csrf.Config          `mapstructure:"csrf"`
}

type App struct {
	Env   string `mapstructure:"env"`
	Debug bool   `mapstructure:"debug"`
}

type Grpc struct {
	AuthServicePort string `mapstructure:"authServicePort"`
	UserServicePort string `mapstructure:"userServicePort"`
}

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
