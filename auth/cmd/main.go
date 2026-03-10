package main

import (
	"auth/internal/config"
	hand "auth/internal/handler"
	"auth/internal/jwtutil"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	logrus.Info(fmt.Sprintf("%+v", cfg))

	if err := jwtutil.LoadKeys(cfg.PrivateKeyPath, cfg.PublicKeyPath); err != nil {
		logrus.Fatal(err)
	}

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		logrus.Fatal(err)
	}
	repo := repository.NewRepository(db)
	srvc := service.NewService(repo)
	handler := hand.NewHandler(srvc)
	srvr := server.NewServer()

	go func() {
		if err := srvr.Run(cfg, handler); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Error starting server: %v", err)
		}
	}()

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srvr.Shutdown(ctx); err != nil {
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
