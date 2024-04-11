package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/RafalSalwa/interview-app-srv/pkg/email"
	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string          `mapstructure:"serviceName"`
	AMQP        rabbitmq.Config `mapstructure:"rabbitmq"`
	Email       email.Config    `mapstructure:"email"`
}

func InitConfig() (*Config, error) {
	cfg := &Config{}
	path, err := getEnvPath()
	if err != nil {
		return nil, err
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}
	return cfg, nil
}

func getEnvPath() (string, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "os.Getwd")
	}
	configPath := ""
	if strings.Contains(getwd, "consumer_service") {
		configPath = fmt.Sprintf("%s/config.yaml", getwd)
	} else {
		configPath = fmt.Sprintf("%s/cmd/consumer_service/config/config.yaml", getwd)
	}
	return configPath, nil
}
