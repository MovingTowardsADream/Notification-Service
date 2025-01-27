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

type Manager interface {
	Close() error
}

type Attr struct {
	Key   string
	Value any
}

func AnyAttr(key string, value any) Attr {
	return Attr{Key: key, Value: value}
}

type HandlingLogger interface {
	Err(err error) Attr
}

type DefaultLogger interface {
	Info(msg string, attrs ...Attr)
	Error(msg string, attrs ...Attr)
	Debug(msg string, attrs ...Attr)
	Warn(msg string, attrs ...Attr)
}

type Logger interface {
	DefaultLogger
	HandlingLogger
	Manager
}

type LogData struct {
	log     *slog.Logger
	logFile *os.File
}

func Setup(env string, filePath *string) (Logger, error) {
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
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case _envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return nil, ErrUnknownEnv
	}

	return &LogData{
		log:     log,
		logFile: file,
	}, nil
}

func (l *LogData) Info(msg string, attrs ...Attr) {
	l.log.Info(msg, convertAttrs(attrs)...)
}

func (l *LogData) Error(msg string, attrs ...Attr) {
	l.log.Error(msg, convertAttrs(attrs)...)
}

func (l *LogData) Debug(msg string, attrs ...Attr) {
	l.log.Debug(msg, convertAttrs(attrs)...)
}

func (l *LogData) Warn(msg string, attrs ...Attr) {
	l.log.Warn(msg, convertAttrs(attrs)...)
}

func convertAttrs(attrs []Attr) []any {
	slogAttrs := make([]any, len(attrs))
	for i, attr := range attrs {
		switch v := attr.Value.(type) {
		case time.Duration:
			slogAttrs[i] = slog.Any(attr.Key, slog.StringValue(v.String()))
		case string:
			slogAttrs[i] = slog.Any(attr.Key, slog.StringValue(v))
		case int, int64, int32, float32, float64:
			slogAttrs[i] = slog.Any(attr.Key, v)
		default:
			slogAttrs[i] = slog.Any(attr.Key, v)
		}
	}
	return slogAttrs
}

func (l *LogData) Err(err error) Attr {
	var msgErr string

	if err != nil {
		msgErr = err.Error()
	}

	return Attr{
		Key:   "error",
		Value: msgErr,
	}
}

func (l *LogData) Close() error {
	if l.logFile != nil {
		return l.logFile.Close()
	}
	return nil
}
