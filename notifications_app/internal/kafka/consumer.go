package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Handler func(message *sarama.ConsumerMessage) error

type Consumer struct {
	consumer sarama.Consumer
	handlers map[string]Handler
}

func NewConsumer(brokers []string) (*Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		handlers: make(map[string]Handler),
	}, nil
}

func (c *Consumer) RegisterHandler(topic string, handler Handler) {
	c.handlers[topic] = handler
}

func (c *Consumer) Start() {
	for topic, handler := range c.handlers {
		partitions, _ := c.consumer.Partitions(topic)
		for _, partition := range partitions {
			pc, _ := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)

			go func(pc sarama.PartitionConsumer) {
				for msg := range pc.Messages() {
					if err := handler(msg); err != nil {
						log.Printf("Ошибка обработки: %v", err)
					}
				}
			}(pc)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
