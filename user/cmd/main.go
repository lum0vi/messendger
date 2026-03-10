package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/config"
	"user/internal/handler"
	"user/internal/repository"
	"user/internal/server"
	"user/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	switch cfg.AppEnv {
	case "local":
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logrus.SetLevel(logrus.DebugLevel)
	case "production":
		gin.SetMode(gin.ReleaseMode)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Infof("%+v", cfg)

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)
	srv := server.NewServer()
	go func() {
		if err := srv.Run(cfg, h.InitRouter()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Error starting server: %v", err)
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown error: %v", err)
	} else {
		logrus.Info("Server shutdown complete.")
	}

	if err := db.DB.Close(); err != nil {
		logrus.Errorf("Database close error: %v", err)
	} else {
		logrus.Info("Database connection closed.")
	}

}
