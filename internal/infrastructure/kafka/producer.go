package kafka

import (
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

func NewProducer(brokers []string) (Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return Producer{}, errors.Wrap(err, "error with sync kafka-producer")
	}

	producer := Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}

	return producer, nil
}

func (p Producer) SendMessage(msg *sarama.ProducerMessage) error {
	_, _, err := p.syncProducer.SendMessage(msg)
	return err
}

func (p Producer) SendMessages(msg []*sarama.ProducerMessage) error {
	return p.syncProducer.SendMessages(msg)
}

func (p Producer) Close() error {
	return p.syncProducer.Close()
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	cfg.Producer.RequiredAcks = sarama.WaitForAll

	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1

	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true

	return sarama.NewSyncProducer(brokers, cfg)
}
