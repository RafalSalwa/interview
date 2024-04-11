//go:build unit

package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArgon2ID(t *testing.T) {
	h, err := Argon2ID("password")
	if err != nil {
		t.Error("should not return error")
	}
	assert.Contains(t, h, "$argon2id")
}

func TestArgon2IDComparePasswordAndHash(t *testing.T) {
	p := "password"
	h := "$argon2id$v=19$m=65536,t=4,p=4$rowcIlsB499gUPbat+2Aow$IJxppAHk5yWKpRemzd5YdYFV9UrLRVG4M+5owhXUrh4"
	_, err := Argon2IDComparePasswordAndHash(p, h)
	if err != nil {
		t.Errorf("compare err")
		return
	}
	assert.NoError(t, err)
}
