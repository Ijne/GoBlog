package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.AsyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case success, ok := <-producer.Successes():
				if !ok {
					return
				}
				log.Printf("Отправлено: topic=%s, partition=%d", success.Topic, success.Partition)
			case err := <-producer.Errors():
				log.Printf("Ошибка: %v", err)
			}
		}
	}()

	return &Producer{producer: producer}, nil
}

func (p *Producer) Send(topic string, key, value []byte) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
