//go:build unit

package config

import (
	"github.com/RafalSalwa/interview-app-srv/pkg/env"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	path, err := env.GetConfigPath("auth_service")
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "auth_service")

	c, err := InitConfig()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.App)
	assert.NotEmpty(t, c.Logger)
	assert.NotEmpty(t, c.Jaeger)
	assert.NotEmpty(t, c.ServiceName)
	assert.NotEmpty(t, c.Email)
	assert.NotEmpty(t, c.JWTToken)
	assert.NotEmpty(t, c.Mongo)
	assert.NotEmpty(t, c.Rabbit)
	assert.NotEmpty(t, c.Redis)
	assert.NotEmpty(t, c.GRPC)
	assert.NotEmpty(t, c.MySQL)
	assert.Equal(t, "Auth Service", c.ServiceName)
}
