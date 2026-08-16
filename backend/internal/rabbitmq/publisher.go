//go:generate mockery --name Publisher --output ../../mocks --outpkg mocks --case underscore
package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/I000000/recly/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type TaskMessage struct {
	TaskID          string             `json:"task_id"`
	UserID          string             `json:"user_id"`
	SelectedIDs     []string           `json:"selected_ids"`
	SelectedWeights map[string]float64 `json:"selected_weights,omitempty"`
	ExcludeIDs      []string           `json:"exclude_ids,omitempty"`
	Direction       string             `json:"direction"`
	Weights         map[string]float64 `json:"weights"`
	Contextual      bool               `json:"contextual"`
}

type Publisher interface {
	PublishRecommendationTask(ctx context.Context, msg TaskMessage) error
	Close() error
}

type AMQPPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewAMQPPublisher(url string) (*AMQPPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}
	q, err := ch.QueueDeclare(
		"recommendation_tasks", // имя
		true,                   // durable
		false,                  // autoDelete
		false,                  // exclusive
		false,                  // noWait
		nil,                    // args
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("queue declare: %w", err)
	}
	return &AMQPPublisher{conn: conn, channel: ch, queue: q}, nil
}

func (p *AMQPPublisher) PublishRecommendationTask(ctx context.Context, msg TaskMessage) error {
	log := logger.FromContext(ctx)
	body, err := json.Marshal(msg)
	if err != nil {
		log.Error("failed to marshal task message", zap.Error(err), zap.String("task_id", msg.TaskID))
		return fmt.Errorf("marshal: %w", err)
	}
	err = p.channel.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		log.Error("failed to publish task to RabbitMQ", zap.Error(err), zap.String("task_id", msg.TaskID))
		return err
	}
	log.Debug("task published to RabbitMQ", zap.String("task_id", msg.TaskID))
	return nil
}

func (p *AMQPPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}
