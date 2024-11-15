package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
)

type LogLevel string

const (
	Debug LogLevel = "DEBUG"
	Info  LogLevel = "INFO"
	Warn  LogLevel = "WARN"
	Error LogLevel = "ERROR"
)

func (l *LogLevel) SlogLevel() slog.Level {
	switch *l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Logger struct {
	logLevel      LogLevel
	defaultLogger *slog.Logger
	fileLogger    *slog.Logger
}

func New(level LogLevel) *Logger {
	return &Logger{
		logLevel:      level,
		defaultLogger: defaultLogger(level.SlogLevel()),
		fileLogger:    fileLogger(level.SlogLevel()),
	}
}

func (l Logger) Warn(msg string) {
	l.defaultLogger.Warn(msg)
}

func (l Logger) Info(msg string) {
	l.defaultLogger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.defaultLogger.Error(msg)
}

func (l Logger) File(msg string) {
	l.fileLogger.Info(msg)
}

func defaultLogger(level slog.Level) *slog.Logger {
	logConfig := &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: nil,
	}
	logHandler := slog.NewJSONHandler(os.Stderr, logConfig)

	return slog.New(logHandler)
}

func fileLogger(level slog.Level) *slog.Logger {
	file, err := os.OpenFile("main.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		panic(fmt.Errorf("create log file: %w", err))
	}

	writer := io.MultiWriter(file, os.Stderr)

	logConfig := &slog.HandlerOptions{
		AddSource:   false,
		Level:       level,
		ReplaceAttr: nil,
	}
	logHandler := slog.NewJSONHandler(writer, logConfig)

	return slog.New(logHandler)
}
