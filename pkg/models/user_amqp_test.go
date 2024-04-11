package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserDBModel_AMQP(t *testing.T) {
	um := UserDBModel{
		Id:               1,
		Username:         "test",
		Email:            "test@test.com",
		VerificationCode: "abcdef",
	}
	m := um.AMQP()
	assert.NotEmpty(t, m)
}
