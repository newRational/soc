package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/newRational/soc/infrastructure/logger"
	"github.com/newRational/soc/internal/api/server"
	"github.com/newRational/soc/internal/model"
	"github.com/newRational/soc/internal/receiver"
)

func runServer(ctx context.Context) error {
	var (
		srv = &http.Server{
			Addr: port,
		}
		wg = sync.WaitGroup{}
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//log.Fatal(err)
			logger.Fatal(err)
		}
	}()

	msgChan := make(chan model.Message, 1000)
	wg.Add(1)
	go func() {
		defer wg.Done()
		receiver.ConsumerGroupLogging(brokers, "responses", msgChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Distribute(msgChan)
	}()

	//log.Printf("listening on %s", port)
	logger.Infof("Listening on %s", port)
	<-ctx.Done()

	wg.Wait()

	//log.Println("shutting down server gracefully")
	logger.Info("Shutting down server gracefully")

	shutdownCtx := context.Context(context.Background())

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	return nil
}
