// //go:build unit

package generator

import (
	"errors"
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerificationCode(t *testing.T) {
	tests := []struct {
		len  int
		name string
	}{
		{
			len:  0,
			name: "len to small (0)",
		},
		{
			len:  6,
			name: "len good (6)",
		},
		{
			len:  20,
			name: "len too big (20)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := RandomString(tt.len)
			if err != nil {
				if errors.Is(err, ErrLengthInvalid) {
					assert.EqualError(t, err, ErrLengthInvalid.Error())
					return
				}
				if errors.Is(err, io.ErrUnexpectedEOF) {
					assert.EqualError(t, err, io.ErrUnexpectedEOF.Error())
					return
				}

				t.Error("unexpected err:", err)
			}
			if code == "" {
				t.Error("empty code val")
			}
			letters := regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
			if !letters(code) {
				t.Error("verification code can contain only letters")
				return
			}
		})
	}
}
