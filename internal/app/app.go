package app

import (
	"fmt"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/handler"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/mailer"
	"github.com/ARUMANDESU/uniclubs-notification-service/internal/rabbitmq"
	"github.com/ARUMANDESU/uniclubs-notification-service/pkg/logger"
	"log/slog"
)

type App struct {
	Rmq      *rabbitmq.Rabbitmq
	handlers *handler.Handlers
	log      *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *App {
	rmq, err := rabbitmq.New(cfg.Rabbitmq, log)
	if err != nil {
		panic(err)
	}

	m := mailer.New(
		cfg.Mailer.Host,
		cfg.Mailer.Port,
		cfg.Mailer.Username,
		cfg.Mailer.Password,
		cfg.Mailer.Sender,
	)

	h := handler.New(m)

	return &App{Rmq: rmq, handlers: h, log: log}
}

func (a *App) Start() error {
	const op = "App.Start"
	log := a.log.With(slog.String("op", op))

	//todo: instead of this handler write proper one
	err := a.Rmq.Consume("user", "user.registered", a.handlers.EmailVerification)
	if err != nil {
		log.Error("consume error", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
