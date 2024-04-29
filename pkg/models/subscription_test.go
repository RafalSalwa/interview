package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionDBModel_TableName(t *testing.T) {
	p := SubscriptionDBModel{}
	assert.Equal(t, "subscription", p.TableName())
}
