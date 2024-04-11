package models

import (
	"encoding/json"
	"fmt"

	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserEvent struct {
	ID               int64
	Username         string
	Email            string
	VerificationCode string
}

func (um *UserDBModel) AMQP() amqp.Publishing {
	ue := UserEvent{
		ID:               um.Id,
		Username:         um.Username,
		Email:            um.Email,
		VerificationCode: um.VerificationCode,
	}
	data, _ := json.Marshal(&ue)
	event := rabbitmq.Event{
		Name:       "customer_account_confirmation_requested",
		ID:         "",
		SequenceID: "",
		Content:    string(data),
	}
	body, _ := json.Marshal(event)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	return message
}
func (u UserEvent) String() string {
	return fmt.Sprintf("%d %s %s %s", u.ID, u.Username, u.Email, u.VerificationCode)
}
