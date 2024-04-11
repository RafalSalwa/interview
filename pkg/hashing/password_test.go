//go:build unit

package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	type args struct {
		plaintext string
		hash      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "hash error",
			args: args{
				plaintext: "interview_test_pass",
				hash:      "$2a$13$yJTkmt1dc184JP/A1uQ3KuSUCbIWS36EiPHpi1xrsmaYRp.9zJavq"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp, err := HashPassword(tt.args.plaintext)
			if err != nil {
				t.Errorf("hash() error = %v", err)
				return
			}
			assert.NotEmpty(t, hp)
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		plaintext string
		hash      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check password",
			args: args{
				plaintext: "interview_test_pass",
				hash:      "$2a$13$yJTkmt1dc184JP/A1uQ3KuSUCbIWS36EiPHpi1xrsmaYRp.9zJavq"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := CheckPasswordHash(tt.args.plaintext, tt.args.hash)
			if !ok {
				t.Errorf("hashed pass = %v, want %v", tt.args.hash, tt.args.plaintext)
				return
			}
		})
	}
}

func TestValidate(t *testing.T) {
	type args struct {
		password    string
		passwordCon string
		errMsg      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check password #1",
			args: args{
				password:    "password",
				passwordCon: "password",
				errMsg:      "Password should contain at least one upper case character"},
		},
		{
			name: "check different password",
			args: args{
				password:    "password",
				passwordCon: "password2",
				errMsg:      "Passwords are not the same"},
		},
		{
			name: "check password len",
			args: args{
				password:    "pass",
				passwordCon: "pass",
				errMsg:      "Password should be between 8 and 16 characters in length"},
		},
		{
			name: "check password lower case (any)",
			args: args{
				password:    "PASSWORD2",
				passwordCon: "PASSWORD2",
				errMsg:      "Password should contain at least one lower case character"},
		},
		{
			name: "check password at least one digit",
			args: args{
				password:    "PASSWORDd",
				passwordCon: "PASSWORDd",
				errMsg:      "Password should contain at least one digit"},
		},
		{
			name: "check password contain special char",
			args: args{
				password:    "PASSWORDd1",
				passwordCon: "PASSWORDd1",
				errMsg:      "Password should contain at least one special character"},
		},
		{
			name: "check good Pass",
			args: args{
				password:    "VeryG00dPass!",
				passwordCon: "VeryG00dPass!",
				errMsg:      ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.password, tt.args.passwordCon)
			if err != nil {
				assert.EqualError(t, err, tt.args.errMsg)
			}
		})
	}
}
