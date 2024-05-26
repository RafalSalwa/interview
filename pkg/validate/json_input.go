package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrEQField struct {
}

func (e ErrEQField) Error() string {
	panic("implement me")
}

func UserInput(r *http.Request, req interface{}) error {
	reqValidator := validator.New()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("cannot decode request data")
	}

	if err := reqValidator.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldErr := range validationErrors {
				return handleFieldError(fieldErr)
			}
		}
	}
	return nil
}

func handleFieldError(fieldErr validator.FieldError) error {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Errorf("field '%s' is required", fieldErr.Field())
	case "eqfield":
		return fmt.Errorf("different values for '%s', but should be same", fieldErr.Param())
	default:
		return fmt.Errorf("field '%s' failed validation for '%s'", fieldErr.Field(), fieldErr.Tag())
	}
}
