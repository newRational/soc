package kafka

import (
	"encoding/json"

	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/model"

	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	ready   chan bool
	msgChan chan<- model.Message
}

func NewConsumerGroup(msgChan chan<- model.Message) ConsumerGroup {
	return ConsumerGroup{
		ready:   make(chan bool),
		msgChan: msgChan,
	}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():

			rm := model.Message{}
			err := json.Unmarshal(message.Value, &rm)
			if err != nil {
				//fmt.Println("Consumer group error", err)
				logger.Error("Consumer group error", err)
			}

			consumer.msgChan <- rm

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			close(consumer.msgChan)
			return nil
		}
	}
}
