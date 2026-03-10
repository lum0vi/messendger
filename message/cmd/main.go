package main

import (
	"context"
	"message/internal/config"
	"message/internal/handler"
	"message/internal/kafka"
	"message/internal/repository"
	"message/internal/server"
	"message/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("%+v", cfg)
	db, err := repository.NewPostgres(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Successfully connected to postgres")
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	hand := handler.NewHandler(svc)
	srv := server.NewServer()
	var addrKafka []string
	addrKafka = append(addrKafka, cfg.KafkaHost+":"+cfg.KafkaPort)

	prod := kafka.NewProducerKafka(addrKafka)
	cons := kafka.NewConsumerKafka(addrKafka, repo, prod)

	go func() {
		cons.Consumer.Start(ctx)
	}()
	go func() {
		if err := srv.Run(hand.InitRouter(), cfg); err != nil && err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()
	logrus.Info("Successfully started server")

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait
	if err := cons.Consumer.Close(); err != nil {
		logrus.Errorf("failed to close consumer: %v", err)
	}
	if err := prod.Producer.Close(); err != nil {
		logrus.Errorf("failed to close producer: %v", err)
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("failed to close DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("server shutdown failed: %v", err)
	}
	select {
	case <-ctx.Done():
		logrus.Info("Shutting down")
		os.Exit(0)
	case <-time.After(2 * time.Second):
		logrus.Info("Timed out waiting for server to shutdown")
		os.Exit(1)
	}
}
