package app

import (
	"encoding/json"
	"fmt"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/domain"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type App struct {
	Rmq *rabbitmq.Rabbitmq
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	rmq, err := rabbitmq.New(cfg.Rabbitmq, log)
	if err != nil {
		panic(err)
	}

	return &App{Rmq: rmq, log: log}
}

func (a *App) Start() error {
	const op = "App.Start"
	log := a.log.With(slog.String("op", op))

	//todo: instead of this handler write proper one
	err := a.Rmq.Consume("notification_queue", func(m amqp091.Delivery) error {
		var msg domain.UserRegister
		err := json.Unmarshal(m.Body, &msg)
		log.Info(msg.Email)
		log.Info(msg.FirstName)
		log.Info(msg.LastName)
		log.Info(msg.Token)
		if err != nil {
			log.Error("failed to unmarshal message: %v", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
