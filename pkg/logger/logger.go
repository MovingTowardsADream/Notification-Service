package logger

import (
	"log/slog"
	"os"
	"time"

	multihandler "Notification_Service/pkg/logger/multi_handler"
)

const (
	_envDev     = "dev"
	_envTesting = "testing"
	_envProd    = "prod"
)

type Logger struct {
	*slog.Logger
	logFile *os.File
}

func Setup(env string, filePath *string) (*Logger, error) {
	var log *slog.Logger
	var file *os.File

	switch env {
	case _envDev:
		if filePath == nil {
			return nil, ErrEmptyFilePath
		}
		file, err := os.OpenFile(*filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			return nil, ErrOpenFile
		}

		log = slog.New(
			multihandler.NewMultiHandler(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
				slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelDebug}),
			),
		)
	case _envTesting:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case _envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return nil, ErrUnknownEnv
	}

	return &Logger{
		Logger:  log,
		logFile: file,
	}, nil
}

func (l *Logger) Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func NewStrArgs(name, value string) slog.Attr {
	return slog.Attr{
		Key:   name,
		Value: slog.StringValue(value),
	}
}

func NewIntArgs(name string, value int) slog.Attr {
	return slog.Attr{
		Key:   name,
		Value: slog.IntValue(value),
	}
}

func NewDurationArgs(name string, value time.Duration) slog.Attr {
	return slog.Attr{
		Key:   name,
		Value: slog.StringValue(value.String()),
	}
}

func (l *Logger) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
