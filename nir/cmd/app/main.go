package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/newRational/soc/infrastructure/kafka"
	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/app/db"
	"github.com/newRational/soc/internal/app/repository"
	"github.com/newRational/soc/internal/app/server"
	"github.com/newRational/soc/internal/app/service"
	"github.com/newRational/soc/internal/model"
	"github.com/newRational/soc/internal/receiver"
	"github.com/newRational/soc/internal/sender"
)

const (
	envFile = "../../prod.env"
)

var brokers = []string{
	"kafka1:9091",
	"kafka2:9092",
	"kafka3:9093",
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	config, err := db.LoadEnv(envFile)
	if err != nil {
		//log.Fatal(err)
		logger.Fatal(err)
	}
	logger.Debug("Load config from environmental")

	database, err := db.NewDb(ctx, config)
	if err != nil {
		//log.Fatal(err)
		logger.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	ordersRepo := repository.NewOrders(database)

	serv := service.NewOrdersService(ordersRepo)

	msgChan := make(chan model.Message, 1000)

	kafkaProducer, err := kafka.NewProducer(brokers, kafka.WithMaxOpenRequests(1), kafka.WithRandomPartitioner(), kafka.WaitForAll(),
		kafka.ReturnSuccesses(true), kafka.ReturnErrors(true), kafka.Idempotent(true), kafka.WithCompressionLevelDefault(), kafka.WithCompressionGZIP())
	if err != nil {
		//fmt.Println(err)
		logger.Error(err)
	}
	sender := sender.NewKafkaSender(kafkaProducer, "responses")

	implemetation := server.NewServer(serv, sender, msgChan)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		receiver.ConsumerGroupLogging(brokers, "requests", msgChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		implemetation.Listen(ctx)
	}()

	wg.Wait()
}
