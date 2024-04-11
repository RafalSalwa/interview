//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	path, err := getEnvPath()
	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "consumer_service")

	c, err := InitConfig()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.ServiceName)
	assert.NotEmpty(t, c.Email)
	assert.NotEmpty(t, c.AMQP)
	assert.Equal(t, "consumer_service", c.ServiceName)
}
