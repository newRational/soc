package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/newRational/soc/infrastructure/kafka"
	"github.com/newRational/soc/infrastructure/logger"
	redisCache "github.com/newRational/soc/infrastructure/redis"
	"github.com/newRational/soc/internal/api/routes"
	"github.com/newRational/soc/internal/api/server"
	"github.com/newRational/soc/internal/api/service"
	"github.com/newRational/soc/internal/sender"
	"github.com/redis/go-redis/v9"
)

const (
	port = ":9000"
)

var brokers = []string{
	"kafka1:9091",
	"kafka2:9092",
	"kafka3:9093",
}

func main() {
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	kafkaProducer, err := kafka.NewProducer(brokers, kafka.WithMaxOpenRequests(1), kafka.WithRandomPartitioner(), kafka.WaitForAll(),
		kafka.ReturnSuccesses(true), kafka.ReturnErrors(true), kafka.Idempotent(true), kafka.WithCompressionLevelDefault(), kafka.WithCompressionGZIP())
	if err != nil {
		//mt.Println(err)
		logger.Error(err)
	}
	sender := sender.NewKafkaSender(kafkaProducer, "requests")

	serv := service.NewService(sender)

	client := redisCache.NewRedis(&redis.Options{
		Addr:     "redis:6379",
		Password: "password",
		DB:       0,
	})

	implemetation := server.NewServer(serv, client)

	http.Handle("/", routes.CreateRouter(implemetation))

	if err := runServer(ctx); err != nil {
		//log.Fatal(err)
		logger.Fatal(err)
	}

}
