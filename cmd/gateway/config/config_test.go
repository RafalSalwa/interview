//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	c, err := InitConfig()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.HTTP)
	assert.NotEmpty(t, c.Grpc)
	assert.NotEmpty(t, c.Auth)
	assert.NotEmpty(t, c.App)
	assert.NotEmpty(t, c.Logger)
	assert.NotEmpty(t, c.Jaeger)
	assert.NotEmpty(t, c.ServiceName)
	assert.Equal(t, c.ServiceName, "api_gateway_service")
}
