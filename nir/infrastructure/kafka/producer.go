package kafka

import (
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

type ConfigOption func(*sarama.Config)

func NewProducer(brokers []string, opts ...ConfigOption) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	producer := &Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}

	return producer, nil
}

func newSyncProducer(brokers []string, opts ...ConfigOption) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()
	for _, o := range opts {
		o(syncProducerConfig)
	}

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	return syncProducer, nil
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	return nil
}

func WithMaxOpenRequests(n int) ConfigOption {
	return func(s *sarama.Config) {
		s.Net.MaxOpenRequests = n
	}
}

func WithRandomPartitioner() ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.Partitioner = sarama.NewRandomPartitioner
	}
}

func WaitForAll() ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.RequiredAcks = sarama.WaitForAll
	}
}

func ReturnSuccesses(n bool) ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.Return.Successes = n
	}
}

func ReturnErrors(n bool) ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.Return.Errors = n
	}
}

func Idempotent(n bool) ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.Idempotent = n
	}
}

func WithCompressionLevelDefault() ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.CompressionLevel = sarama.CompressionLevelDefault
	}
}

func WithCompressionGZIP() ConfigOption {
	return func(s *sarama.Config) {
		s.Producer.Compression = sarama.CompressionGZIP
	}
}
