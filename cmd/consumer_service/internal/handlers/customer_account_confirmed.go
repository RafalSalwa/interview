package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
)

type CustomerAccountActivatedEventEmail struct {
	UserID           int    `json:"customer_id"`
	TrackingClientID string `mapstructure:"tracking_client_id"`
	Country          string `mapstructure:"country_code"`
}

func WrapHandleCustomerAccountConfirmedEmail(event rabbitmq.Event) error {
	var data CustomerAccountActivatedEventEmail
	err := json.Unmarshal([]byte(event.Content), &data)

	if err != nil {
		return err
	}

	return HandleCustomerAccountConfirmEmail(data)
}

func HandleCustomerAccountConfirmEmail(payload CustomerAccountActivatedEventEmail) error {
	fmt.Println("HandleCustomerAccountConfirmEmail", payload)
	return nil
}
