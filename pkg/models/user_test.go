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
	err := dbu.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.Equal(t, false, dbu.Active)
	assert.Equal(t, "user", dbu.TableName())
}
