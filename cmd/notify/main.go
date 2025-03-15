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
	ctx := context.Background()

	cfg := config.MustLoad()

	log, err := logger.Setup(cfg.Log.Level, cfg.Log.Path)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := log.Close(); err != nil {
			log.Error("failed to close logger", log.Err(err))
		}
	}()

	application := app.New(ctx, log, cfg)

	go func() {
		if errOutbox := application.OutboxWorker.WorkerRun(); errOutbox != nil {
			log.Error("outboxWorker error: ", log.Err(errOutbox))
		}
	}()

	go func() {
		if errServ := application.MetricsServer.Run(); errServ != nil {
			log.Error("server with metrics was shut down due to an error: ", log.Err(errServ))
		}
	}()

	go func() {
		application.MessagingServer.MustRun()
	}()

	go func() {
		if errServ := application.Server.Run(); errServ != nil {
			log.Error("server was shut down due to an error: ", log.Err(errServ))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
	case <-application.MessagingServer.Notify():
	}

	log.Info("starting graceful shutdown")

	application.Server.Shutdown()

	if err := application.MessagingServer.Shutdown(); err != nil {
		log.Error("messaging shutdown error", log.Err(err))
	}

	// TODO outbox stop

	if err := application.Tracer.Close(ctx); err != nil {
		log.Error("tracer shutdown error", log.Err(err))
	}

	application.Storage.Close()

	if err := application.MetricsServer.Shutdown(ctx); err != nil {
		log.Error("metrics server shutdown error", log.Err(err))
	}

	log.Info("gracefully stopped")
}
