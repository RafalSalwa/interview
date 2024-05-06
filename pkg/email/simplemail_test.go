//go:build unit

package email

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient(t *testing.T) {
	config := Config{
		Host:        "mailpit",
		Port:        1025,
		From:        "interview@example.com",
		TemplateDir: "../../templates",
	}
	assert.NotEmpty(t, NewClient(config))
}

func TestSendVerificationEmail(t *testing.T) {
	config := Config{
		Host:        "mailpit",
		Port:        1025,
		From:        "interview@example.com",
		TemplateDir: "../../templates",
	}
	mailer := NewClient(config)
	mail := UserEmailData{
		Username:         "test@test.com",
		Email:            "test@test.com",
		VerificationCode: "asdf",
	}
	mailer.SendVerificationEmail(mail)
}
