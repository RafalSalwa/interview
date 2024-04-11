package hashing

import (
	"regexp"

	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

type MismatchError struct {
	error string
}

func (m MismatchError) Error() string {
	return "boom"
}

type ValidationError struct {
	Message string
	Field   string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Validate(password string, PasswordConfirm string) error {
	if password != PasswordConfirm {
		return &ValidationError{
			Message: "Passwords are not the same",
			Field:   "passwordConfirm",
		}
	}

	if len(password) < 8 || len(password) > 32 {
		return &ValidationError{
			Message: "Password should be between 8 and 32 characters in length",
			Field:   "password",
		}
	}

	done, err := regexp.MatchString("([a-z])+", password)
	if err != nil {
		return err
	}
	if !done {
		return &ValidationError{
			Message: "Password should contain at least one lower case character",
			Field:   "password",
		}
	}

	done, err = regexp.MatchString("([A-Z])+", password)
	if err != nil {
		return err
	}
	if !done {
		return &ValidationError{
			Message: "Password should contain at least one upper case character",
			Field:   "password",
		}
	}

	done, err = regexp.MatchString("([0-9])+", password)
	if err != nil {
		return err
	}
	if !done {
		return &ValidationError{
			Message: "Password should contain at least one digit",
			Field:   "password",
		}
	}

	done, err = regexp.MatchString("([!@#$%^&*.?-])+", password)
	if err != nil {
		return err
	}
	if !done {
		return &ValidationError{
			Message: "Password should contain at least one special character",
			Field:   "password",
		}
	}

	err = passwordvalidator.Validate(password, 70)
	if err != nil {
		return err
	}
	return nil
}
