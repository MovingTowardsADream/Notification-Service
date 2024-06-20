# Notification Service

A notification service is a software system that manages the sending of notifications across different channels and devices.

<p align="center">
  <img src="https://github.com/MovingTowardsADream/Notification-Service/blob/main/architecture%20design/Notification-Service.jpg" alt="Notification Service Architecture Design" width="600"/>
</p>

## The Launch of
To run an application in a container, you need to configure the command:
```bash
make build
```
If you are launching the application for the first time, do not forget to make the migration:
```bash
make migration-new-db
```

## Environment Variables and Configuration
For the application to work correctly, you must specify environment variables in the .env file in the root directory. Below are the variables themselves and a brief description:
`GIN_MODE` - set to debug mode if debugging is necessary.

`POSTGRES_USER`, `POSTGRES_DB`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_PASSWORD` - parameters for initializing the Postgresql database in docker-compose.

`PG_URL` - link to connect to Postgresql.

`RMQ_URL` - link to the rabbitmq queue.

There is also a config.yaml file in which the remaining data is specified.
