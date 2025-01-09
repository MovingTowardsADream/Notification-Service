package logger

import (
	"log/slog"
	"os"

	"Notification_Service/pkg/logger/multi_handler"
)

const (
	_envDev     = "dev"
	_envTesting = "testing"
	_envProd    = "prod"
)

type ExtendedLoggerParams interface {
	With(args ...any) *slog.Logger
}

type DefaultLoggerFunc interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type DefaultLogger interface {
	DefaultLoggerFunc
	ExtendedLoggerParams
}

type Logger struct {
	DefaultLogger
	logFile *os.File
}

func (l *Logger) Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Setup(env string, filePath *string) (*Logger, error) {
	var log *slog.Logger
	var file *os.File

	switch env {
	case _envDev:
		if filePath == nil {
			return nil, ErrEmptyFilePath
		}
		file, err := os.OpenFile(*filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, ErrOpenFile
		}

		log = slog.New(
			multi_handler.NewMultiHandler(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
				slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		)
	case _envTesting:
		log = slog.New(
			multi_handler.NewMultiHandler(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			),
		)
	case _envProd:
		log = slog.New(
			multi_handler.NewMultiHandler(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
			),
		)
	default:
		return nil, ErrUnknownEnv
	}

	return &Logger{
		DefaultLogger: log,
		logFile:       file,
	}, nil
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
