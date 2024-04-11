package validate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func UserInput(r *http.Request, req interface{}) error {
	reqValidator := validator.New()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("cannot decode request data")
	}

	if err := reqValidator.Struct(req); err != nil {
		return fmt.Errorf("data validation failed with reason: %s", err.Error())
	}
	return nil
}
