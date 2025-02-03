package gotests

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	// package for initializing the PostgreSQL driver, which is not used directly.
	_ "github.com/lib/pq"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func initDatabaseIntegration(ctx context.Context) func() {
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "notify",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "admin",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(10 * time.Second),
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"5432/tcp": []nat.PortBinding{
						{HostIP: "0.0.0.0", HostPort: "5432"},
					},
				}
			},
		},
		Started: true,
	})

	if err != nil {
		panic(err)
	}

	return func() {
		_ = pgContainer.Terminate(ctx)
	}
}

func initMessagingIntegration(ctx context.Context) func() {
	rmqContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "rabbitmq:3.8.12-management",
			ExposedPorts: []string{"5672/tcp", "15672/tcp"},
			Env: map[string]string{
				"RABBITMQ_DEFAULT_USER": "amqp",
				"RABBITMQ_DEFAULT_PASS": "admin",
			},
			WaitingFor: wait.ForListeningPort("5672/tcp").WithStartupTimeout(10 * time.Second),
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"5672/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "5672"}},
					"15672/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "15672"}},
				}
			},
		},
		Started: true,
	})

	if err != nil {
		panic(err)
	}

	return func() {
		_ = rmqContainer.Terminate(ctx)
	}
}
