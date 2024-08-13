package main

import (
	"Notification_Service/configs"
	"Notification_Service/internal/app"
	"Notification_Service/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init configuration
	cfg := configs.MustLoad()

	fmt.Println(cfg)

	// Init logger
	log := logger.SetupLogger(cfg.Log.Level)

	// Init application
	application := app.New(log, cfg)

	// Run servers
	go func() {
		application.Server.Run()
	}()

	go func() {
		application.RMQServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
	case <-application.RMQServer.Notify():
	}

	log.Info("Starting graceful shutdown")

	if err := application.RMQServer.Shutdown(); err != nil {
		log.Error("RMQServer.Shutdown error", logger.Err(err))
	}

	application.Server.Shutdown()

	application.DB.Close()

	log.Info("Gracefully stopped")
}
