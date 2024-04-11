package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderDBModel_TableName(t *testing.T) {
	o := OrderDBModel{}
	assert.Equal(t, "orders", o.TableName())
}
