package hashing

import (
	"fmt"
	"regexp"

	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type (
	MismatchError struct {
		error string
	}
	ValidationError struct {
		Message string
		Field   string
	}
)

const (
	PassField             = "password"
	PassMinLength         = 8
	PassMaxLength         = 32
	EntropyMinForPassword = 70
	BCryptCost            = 13
)

func (m MismatchError) Error() string {
	return m.error
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BCryptCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Validate(password, passwordConfirm string) error {
	if password != passwordConfirm {
		return &ValidationError{Message: "Passwords are not the same", Field: "passwordConfirm"}
	}

	if len(password) < PassMinLength || len(password) > PassMaxLength {
		return &ValidationError{Message: "Password should be between 8 and 32 characters in length", Field: PassField}
	}

	r, err := regexp.Compile(`^(?=.*[[:lower:]])(?=.*[[:upper:]])(?=.*[[:digit:]]).+$`)
	if err != nil {
		return err
	}

	if !r.MatchString(password) {
		msg := fmt.Sprintf("Password must contain at least: %s %s %s",
			"one lower case character",
			"one upper case character",
			"one number",
		)
		return &ValidationError{Message: msg, Field: PassField}
	}

	done, err := regexp.MatchString("([!@#$%^&*.?-])+", password)
	if err != nil {
		return err
	}
	if !done {
		return &ValidationError{Message: "Password should contain at least one special character", Field: PassField}
	}

	err = passwordvalidator.Validate(password, EntropyMinForPassword)
	if err != nil {
		return err
	}
	return nil
}
