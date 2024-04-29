package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/RafalSalwa/auth-api/pkg/email"
	"github.com/RafalSalwa/auth-api/pkg/rabbitmq"
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
		return nil, fmt.Errorf(" viper read condig %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf(" viper.Unmarshal %w", err)
	}
	return cfg, nil
}

func getEnvPath() (string, error) {
	getwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf(" os.getwd %w", err)
	}
	configPath := ""
	if strings.Contains(getwd, "consumer_service") {
		configPath = fmt.Sprintf("%s/config.yaml", getwd)
	} else {
		configPath = fmt.Sprintf("%s/cmd/consumer_service/config/config.yaml", getwd)
	}
	return configPath, nil
}
