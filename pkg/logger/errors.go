package logger

import (
	"errors"
)

var (
	ErrOpenFile      = errors.New("failed to open file")
	ErrUnknownEnv    = errors.New("unknown environment")
	ErrEmptyFilePath = errors.New("empty file path")
)
