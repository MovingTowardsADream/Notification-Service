package main

import (
	"Notification_Service/configs"
	"Notification_Service/internal/app"
	"Notification_Service/pkg/logger"
)

func main() {
	// Init configuration
	cfg := configs.MustLoad()

	// Init logger
	log := logger.SetupLogger(cfg.Log.Level)

	// Init application
	application := app.New(log, cfg)

	_ = application
}
