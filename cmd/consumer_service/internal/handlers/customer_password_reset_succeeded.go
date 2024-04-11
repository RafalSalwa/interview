package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
)

type CustomerPasswordResetSucceeded struct {
	CustomerID   int    `json:"customer_id"`
	CustomerUUID string `json:"customer_uuid"`
}

func WrapHandleCustomerPasswordResetRequestedSuccessed(event rabbitmq.Event) error {
	var data CustomerPasswordResetSucceeded
	err := json.Unmarshal([]byte(event.Content), &data)

	if err != nil {
		return err
	}

	return CustomerPasswordResetSuccessEmail(data)
}

func CustomerPasswordResetSuccessEmail(payload CustomerPasswordResetSucceeded) error {
	fmt.Println("CustomerPasswordResetSuccessEmail", payload)
	return nil
}
