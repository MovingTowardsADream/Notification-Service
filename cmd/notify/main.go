package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"Notification_Service/internal/application/app"
	"Notification_Service/internal/infrastructure/config"
	"Notification_Service/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log, err := logger.SetupLogger(cfg.Log.Level, cfg.Log.Path)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := log.Close(); err != nil {
			log.Error("failed to close logger", log.Err(err))
		}
	}()

	application := app.New(context.Background(), log, cfg)

	go func() {
		if errServ := application.Server.Run(); errServ != nil {
			log.Error("server was shut down due to an error: ", log.Err(errServ))
		}
	}()

	go func() {
		application.MessagingServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
	case <-application.MessagingServer.Notify():
	}

	log.Info("starting graceful shutdown")

	if err := application.MessagingServer.Shutdown(); err != nil {
		log.Error("messaging shutdown error", log.Err(err))
	}

	application.Server.Shutdown()

	application.Storage.Close()

	log.Info("gracefully stopped")
}
