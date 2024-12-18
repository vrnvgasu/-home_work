package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type App interface {
	SetEventsSent(ctx context.Context, eventIDs []uint64) error
}

type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
	Shutdown() error
}

type Sender struct {
	consumer Consumer
	app      App
	logger   Logger
	done     chan error
}

func NewSender(consumer Consumer, app App, logger Logger) *Sender {
	return &Sender{
		consumer: consumer,
		app:      app,
		logger:   logger,
		done:     make(chan error),
	}
}

func (s *Sender) Run() error {
	deliveries, err := s.consumer.Consume()
	if err != nil {
		s.logger.Info("sender Run Consume: " + err.Error())

		return err
	}

	go s.Handle(deliveries)

	return nil
}

func (s *Sender) Handle(deliveries <-chan amqp.Delivery) {
	ctx := context.Background()
	for d := range deliveries {
		msg := fmt.Sprintf("got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)

		event := &storage.Event{}
		err := json.Unmarshal(d.Body, event)
		if err != nil {
			s.logger.Error("sender Handle Unmarshal: " + err.Error())
		}
		if err = s.app.SetEventsSent(ctx, []uint64{event.ID}); err != nil {
			s.logger.Error("scheduler eventsToSend SetEventsSent: " + err.Error())
		}

		s.logger.Info(msg)
		d.Ack(true)
	}

	s.logger.Info("sender Handle: deliveries channel closed")

	s.done <- nil
}

func (s *Sender) Stop() {
	s.logger.Info("sender Stop")

	if err := s.consumer.Shutdown(); err != nil {
		s.logger.Error("sender Stop Shutdown: %s" + err.Error())
	}

	<-s.done
}
