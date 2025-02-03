package gotests

import (
	"time"

	"Notification_Service/internal/infrastructure/config"
)

func initTestConfig() *config.Config {
	return &config.Config{
		App: config.App{
			Name:         "notification-service",
			Version:      "v1.0.0",
			CountWorkers: 24,
			Timeout:      5 * time.Second,
		},
		GRPC: config.GRPC{
			Port:    8082,
			Timeout: 5 * time.Second,
		},
		Storage: config.Storage{
			PoolMax:      2,
			URL:          "postgres://postgres:admin@localhost:5432/notify?sslmode=disable",
			ConnAttempts: 10,
			ConnTimeout:  1 * time.Second,
		},
		Messaging: config.Messaging{
			Server: config.MessagingServer{
				RPCExchange:     "rpc_server",
				GoroutinesCount: 11,
				WaitTime:        2 * time.Second,
				Attempts:        10,
				Timeout:         2 * time.Second,
			},
			Client: config.MessagingClient{
				RPCExchange: "rpc_client",
				WaitTime:    2 * time.Second,
				Attempts:    10,
				Timeout:     2 * time.Second,
			},
			URL:    "amqp://amqp:admin@localhost:5672/",
			Topics: []string{"mail_notify", "phone_notify"},
		},
		Log: config.Log{
			Level: "testing",
		},
		Security: config.Security{
			PasswordSalt: "salt",
		},
	}
}
