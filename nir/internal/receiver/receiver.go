package receiver

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/newRational/soc/infrastructure/kafka"
	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/model"

	"github.com/IBM/sarama"
)

func ConsumerGroupLogging(brokers []string, topic string, msgChan chan<- model.Message) {
	keepRunning := true
	//log.Println("Starting a new Sarama consumer")
	logger.Info("Starting a new Sarama consumer")

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	config.Consumer.Group.ResetInvalidOffsets = true

	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	config.Consumer.Group.Session.Timeout = 60 * time.Second

	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second

	const BalanceStrategy = "roundrobin"
	switch BalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		//log.Panicf("Unrecognized consumer group partition assignor: %s", BalanceStrategy)
		logger.Fatalf("Unrecognized consumer group partition assignor: %s", BalanceStrategy)
	}

	consumer := kafka.NewConsumerGroup(msgChan)
	group := topic + "-group"

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		//log.Panicf("Error creating consumer group client: %v", err)
		logger.Fatalf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, []string{topic}, &consumer); err != nil {
				//log.Panicf("Error from consumer: %v", err)
				logger.Fatalf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.Ready()
	//log.Println("Sarama consumer up and running!...")
	logger.Info("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			//log.Println("terminating: context cancelled")
			logger.Info("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			//log.Println("terminating: via signal")
			logger.Info("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		//log.Panicf("Error closing client: %v", err)
		logger.Fatalf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		//log.Println("Resuming consumption")
		logger.Info("Resuming consumption")
	} else {
		client.PauseAll()
		//log.Println("Pausing consumption")
		logger.Info("Resuming consumption")
	}

	*isPaused = !*isPaused
}
