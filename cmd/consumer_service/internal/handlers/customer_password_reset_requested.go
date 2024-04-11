package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
)

type CustomerPasswordResetRequested struct {
	CustomerID   int    `json:"customer_id"`
	CustomerUUID string `json:"customer_uuid"`
	ResetCode    string `json:"reset_code"`
}

func WrapHandleCustomerPasswordResetRequestedEmail(event rabbitmq.Event) error {
	var data CustomerPasswordResetRequested
	err := json.Unmarshal([]byte(event.Content), &data)

	if err != nil {
		return err
	}

	return CustomerPasswordResetRequestEmail(data)
}

func CustomerPasswordResetRequestEmail(payload CustomerPasswordResetRequested) error {
	fmt.Println("CustomerPasswordResetRequestEmail", payload)
	return nil
}
