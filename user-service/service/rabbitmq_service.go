package service

import (
	"context"
	"fmt"
	"micro-warehouse/user-service/configs"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type EmailPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
}

type RabbitMQServiceInterface interface {
	PublishEmail(ctx context.Context, payload EmailPayload) error
	Close() error
}

type rabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  configs.Config
}

// Close implements [RabbitMQServiceInterface].
func (r *rabbitMQService) Close() error {
	panic("unimplemented")
}

// PublishEmail implements [RabbitMQServiceInterface].
func (r *rabbitMQService) PublishEmail(ctx context.Context, payload EmailPayload) error {
	panic("unimplemented")
}

func NewRabbitMQService(config configs.Config) (RabbitMQServiceInterface, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQ.Username, config.RabbitMQ.Password, config.RabbitMQ.Host, config.RabbitMQ.Port))
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 1 %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 2 %v", err)
		return nil, err
	}

	return &rabbitMQService{
		conn:    conn,
		channel: ch,
		config:  config,
	}, nil
}
