package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentDBModel_TableName(t *testing.T) {
	p := PaymentDBModel{}
	assert.Equal(t, "payment", p.TableName())
}
