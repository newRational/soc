package sender

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/model"
)

type producer interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
	Close() error
}

type KafkaSender struct {
	producer producer
	topic    string
}

func NewKafkaSender(producer producer, topic string) *KafkaSender {
	return &KafkaSender{
		producer,
		topic,
	}
}

func (s *KafkaSender) SendMessage(message model.Message) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		//fmt.Println("Send message marshal error", err)
		logger.Error("Send message marshal error", err)
		return err
	}

	_, _, err = s.producer.SendSyncMessage(kafkaMsg)

	if err != nil {
		//fmt.Println("Send message connector error", err)
		logger.Error("Send message connector error", err)
		return err
	}

	return nil
}

func (s *KafkaSender) buildMessage(message model.Message) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)

	if err != nil {
		//fmt.Println("Send message marshal error", err)
		logger.Error("Send message marshal error", err)
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(message.ID)),
	}, nil
}
