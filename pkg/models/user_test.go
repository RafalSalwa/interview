package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModels(t *testing.T) {
	dbu := UserDBModel{
		Id: 1,
	}
	assert.NotEmpty(t, dbu)
	assert.False(t, dbu.Active)
	assert.Equal(t, "user", dbu.TableName())
}
