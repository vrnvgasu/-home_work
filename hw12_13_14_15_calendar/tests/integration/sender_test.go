//go:build integration

package integration

import (
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
)

type senderRequest struct {
	ID     uint64 `json:"id" db:"id"`
	IsSent bool   `json:"is_sent" db:"is_sent"`
}

func (s *SuiteIntegrationTest) TestHandleSender() {
	var id uint64
	err := s.db.QueryRow(`insert into event (title, start_at, end_at, description, owner_id, send_before, is_sent) 
		values ('title_old', now(), now(), 'description_old', '11', '100', false) returning id`).Scan(&id)
	s.Require().NoError(err)
	eventRequest := senderRequest{
		ID:     id,
		IsSent: false,
	}
	s.publishToRabbit(eventRequest)

	time.Sleep(1 * time.Second)
	var isSent bool
	err = s.db.QueryRow("select is_sent from event where id = $1", id).Scan(&isSent)
	s.Require().NoError(err)

	s.Require().True(isSent)
}

func (s *SuiteIntegrationTest) publishToRabbit(req senderRequest) {
	connection, err := amqp.Dial(s.amqpURI)
	s.Require().NoError(err)

	channel, err := connection.Channel()
	s.Require().NoError(err)

	err = channel.ExchangeDeclare(
		s.exchangeName,
		s.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	s.Require().NoError(err)

	msg, err := json.Marshal(req)
	s.Require().NoError(err)

	err = channel.Publish(
		s.exchangeName,
		s.routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            msg,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	s.Require().NoError(err)
}
