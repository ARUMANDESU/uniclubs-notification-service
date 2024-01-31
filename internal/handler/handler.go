package handler

import (
	"encoding/json"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/mailer"
	"github.com/rabbitmq/amqp091-go"
)

type Handlers struct {
	mailSender mailer.MailSender
}

func New(sender mailer.MailSender) *Handlers {
	return &Handlers{mailSender: sender}
}

func (h Handlers) EmailVerification(m amqp091.Delivery) error {
	const op = "Handlers.EmailVerification"

	var msg domain.UserRegister

	err := json.Unmarshal(m.Body, &msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = h.mailSender.Send(msg.Email, "user_welcome.tmpl", msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
