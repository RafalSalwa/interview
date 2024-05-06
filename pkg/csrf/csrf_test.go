//go:build unit

package csrf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeToken(t *testing.T) {
	// Test case 1: Empty salt
	cfg := Config{salt: ""}
	token := MakeToken(cfg)
	assert.NotEmpty(t, token)

	// Test case 2: Non-empty salt
	cfg = Config{salt: "my_salt"}
	token = MakeToken(cfg)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	// Test case 1: Valid token
	cfg := Config{salt: "my_salt"}
	token := MakeToken(cfg)
	isValid := ValidateToken(token, cfg)
	if !isValid {
		t.Error("ValidateToken() with valid token failed; expected true, got false")
	}

	// Test case 2: Invalid token
	cfg = Config{salt: "my_salt"}
	isValid = ValidateToken("invalid_token", cfg)
	if isValid {
		t.Error("ValidateToken() with invalid token passed; expected false, got true")
	}
}
