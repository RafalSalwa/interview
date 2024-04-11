package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlan_TableName(t *testing.T) {
	p := Plan{}
	assert.Equal(t, "plan", p.TableName())
}
